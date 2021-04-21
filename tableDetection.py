"""Importing necessary modules"""

import json
import re

import cv2
import numpy as np
import pytesseract as pt
from utils.difFixer import fix_string as fx
from utils.difFixer import rreplace as rep


def get_venue(param):
    """Function to find venue for classes"""
    try:
        term = re.findall(r"[A-z0-9]{2,5}[0-9A-Z]{1,5}[A-Za-z]{0,2}\b", param)
        print("term: ", term)
        if (term[0]).startswith("1") or len(term[0]) == 3:
            term.pop(0)
        # Fix venue errors during OCR runtime
        for i in term:
            i = i.replace("$", "S")
            i = i.replace("Sjt", "SJT")
            i = i.replace("a", "1")
            i = i.replace("Sit", "SJT")
            i = i.replace("SjT", "SJT")
            i = i.replace("SII", "SJT")
            i = i.replace("SJI", "SJT")
            i = i.replace("SyTa", "SJT4")
            i = i.replace("B", "8")
            i = i.replace("SMVGO", "SMVG")
            i = i.replace("SMG", "SMVG")
            i = i.replace("5", "S")
            i = i.replace("SIT", "SJT")
            i = i.replace("I", "T")
            if "S" in i[3:]:
                i = rep(i, "S", "5", i.count("S") - 1)
                venue = i
            elif "7" in i[0:3]:
                i = i.replace("7", "T", 1)
                venue = i

            elif i.startswith("STS"):
                venue = term[1]
            else:
                if len(term) == 3:
                    venue = term[1]
                elif len(term) == 4:
                    venue = term[2]
                else:
                    venue = i
            return venue
    except:
        venue = None
        return venue


def value_error_handler(param):
    """Error handling"""
    print(param)
    return None


def fetch_data(image):
    """Function to fetch timetable from image"""
    # Gel all pixels in the image - where BGR = (51,255,204), OpenCV colors order is BGR not RGB (green color)
    global slot, course_name, course_name_raw
    slot = None
    course_name_raw = None
    course_name = None
    course_type = None
    gray = np.all(image == (51, 255, 204), 2)  # gray is a logical matrix with True
    # Convert logical matrix to uint8
    gray = gray.astype(np.uint8) * 255
    cnts = cv2.findContours(gray, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_NONE)[0]
    boxes = []
    cnts.reverse()
    data = list()

    def write_json(data, filename="data.json"):
        """ Function to write a json file"""
        with open(filename, "w") as f:
            json.dump(data, f, indent=4)

    for c in cnts:
        (x, y, w, h) = cv2.boundingRect(c)
        boxes.append([x, y, x + w, y + h])
        cv2.rectangle(image, (x, y), (x + w, y + h), (0, 0, 255), 2)
        crop = gray[y : y + h, x : x + w]

        cv2.imwrite("dump1.png", crop)
        dump_file = "dump1.png"
        # test_file = ocr(
        # filename=dump_file, overlay=True, language="eng"
        # )
        text = pt.image_to_string(crop)
        # test_file = json.loads(test_file)
        # text = test_file["ParsedResults"][0]["ParsedText"]
        text = re.sub("[‘]", "", text)
        text = re.sub("[\[(){}<>‘\]|/]", "J", text)

        try:
            slot = re.findall(r"^[A-Za-z0-9]{1,3}[0-9A-Za-z]{0,2}\b", text)
            print("SLOT:", slot)
            slot = slot[0]
            slot = fx(slot)
            course_name_raw = re.findall(r"[A-z0-9]{3,6}[0-9]{1,4}", text)
            print(course_name_raw)
            if len(course_name_raw) > 1:
                for i in course_name_raw:
                    if len(i) < 7:
                        course_name_raw.pop()
                print("CNAME: ", course_name_raw)
                if (course_name_raw[0]).startswith("1"):
                    course_name = course_name_raw[1]
                else:
                    course_name = course_name_raw[0]
                    course_name = fx(course_name)
            elif len(course_name_raw) == 0 or course_name_raw is None:
                course_name = None

            # print(course_name)
            course_code = re.findall(r"[ETH,SS,ELA,LO]{2,3}\b", text)

            course_type = "Lab" if course_code[0] in ("ELA", "LO") else "Theory"

            venue = get_venue(text)
            if course_name is not None and venue == course_name:
                fix_venue = course_name_raw[1]
                unwanted = [foo for foo in ["ETH", "ELA", "LO", "FLA", "BLA", "LO-"]]
                for goo in unwanted:
                    fix_venue = fix_venue.lstrip(goo)
                print("boom")
                venue = get_venue(fix_venue)

        except IndexError:
            if len(slot) == 0 or slot is None:
                slot = None
            venue = None
        except UnboundLocalError:
            course_name = None
            # slot=None

        finally:
            # writing data to json
            slot_data = {
                "Parsed_Data": text,
                "Slot": slot,
                "Course_Name": course_name,
                "Course_type": course_type,
                "Venue": venue,
            }
            data.append(slot_data)
        cv2.imwrite("gray.png", gray)

        write_json(data)
        # returning data
    return {"Slots": data}


def fetch_text_timetable(text):
    data, slots = [], []
    slots += re.findall(
        r"[A-Z]{1,3}[0-9]{1,2}[\D]{1}[A-Z]{3}[0-9]{4}[\D]{1}[A-Z]{2,3}[\D]{1}[A-Z]{2,4}[0-9]{3,4}[A-Z]{0,1}[\D]{1}[A-Z]{3}",
        text,
    )
    for single_slot in slots:
        slot = re.findall(r"[A-Z]{1,3}[0-9]{1,2}\b", single_slot)[0]
        course_name = re.findall(r"[A-Z]{3}[0-9]{4}\b", single_slot)[0]
        course_code = re.findall(r"[ETH,SS,ELA,LO]{2,3}\b", single_slot)
        course_type = "Lab" if course_code[0] in ("ELA", "LO") else "Theory"
        venue = re.findall(r"[A-Z]{2,4}[0-9]{3,4}[A-Z]{0,1}\b", single_slot)[1]
        slot_data = {
            "Parsed_Data": single_slot,
            "Slot": slot,
            "Course_Name": course_name,
            "Course_type": course_type,
            "Venue": venue,
        }
        data.append(slot_data)
    return {"Slots": data}
