"""Importing necessary modules"""

import json
import re

import cv2
import numpy as np

from utils.api import ocr_space_file as ocr
from utils.dif_fixer import fix_string as fx


def get_venue(param):
    """Function to find venue for classes"""
    try:
        term = re.findall(
            r"[0-9A-Za-z]{2,5}[0-9]{2,5}[A-Za-z]{0,1}\b", param
        )
        # Fix venue errors during OCR runtime
        for i in term:
            i = i.replace("SMVGO", "SMVG")
            i = i.replace("SMG","SMVG")
            i = i.replace("$", "S")
            i = i.replace("5","S")
            i = i.replace("SIT","SJT")
            i = i.replace("I","T")
            if i.startswith("STS"):
                venue = term[1]
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
    gray = np.all(image == (51, 255, 204), 2)  # gray is a logical matrix with True
    # gray = image

    # Convert logical matrix to uint8
    gray = gray.astype(np.uint8) * 255
    # print(gray)
    cnts = cv2.findContours(gray, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_NONE)[0]
    # print(cnts)
    boxes = []
    cnts.reverse()
    data = list()

    def write_json(data, filename="data.json"):
        """ Function to write a json file"""
        with open(filename, "w") as f:
            json.dump(data, f, indent=4)

    i = 0
    for c in cnts:
        (x, y, w, h) = cv2.boundingRect(c)
        boxes.append([x, y, x + w, y + h])
        cv2.rectangle(image, (x, y), (x + w, y + h), (0, 0, 255), 2)
        crop = gray[y: y + h, x: x + w]

        cv2.imwrite("dump1.png",crop)
        dump_file = ("dump1.png")
        test_file = ocr(
            filename=dump_file, overlay=True, language="eng"
        )
        test_file = json.loads(test_file)
        text = test_file["ParsedResults"][0]["ParsedText"]
        print(text)
        try:
            slot = re.findall(r"^[A-Za-z]{1,3}[0-9A-Za-z]{1,2}\b", text)
            print(slot)
            slot =slot[0]
            slot = fx(slot)
            global course_name_raw
            course_name_raw = re.findall(r"[A-Za-z]{3,6}[0-9]{1,4}", text)
            print(course_name_raw)
            course_name = course_name_raw[0]
            if course_name is None or len(course_name)==0:
                course_name = None
            course_name = fx(course_name)
            print(course_name)
            course_code = re.findall(r"[ETH,SS,ELA,LO]{2,3}\b", text)

            course_type = "Lab" if course_code[0] in ("ELA", "LO") else "Theory"

            venue = get_venue(text)
            if venue == course_name:
                fix_venue = course_name_raw[1]
                unwanted = [foo for foo in ["ETH","ELA","LO"]]
                for goo in unwanted:
                    fix_venue = fix_venue.lstrip(goo)
                venue = get_venue(fix_venue)
                print(course_name_raw[1],fix_venue, venue)

            # print(slot, course_name, course_type,venue)
        except IndexError:
            if len(slot) == 0 or slot== None:
                slot = value_error_handler(slot)
            venue = None

        # writing data to json
        slot_data = {
            "Parsed_Data": text,
            "Slot": slot,
            "Course_Name": course_name,
            #"Course_Code": (course_code[0]).upper(),
            "Course_type": course_type,
            "Venue": venue,
        }
        data.append(slot_data)
    cv2.imwrite("gray.png", gray)

    write_json(data)
    # returning data
    return {"Slots": data}
