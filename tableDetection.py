"""Importing necessary modules"""

import json
import re

import cv2
import numpy as np
import pytesseract as pt


def get_venue(venue):
    """Function to find venue for classes"""
    term = re.findall(
        r"[TT,SJT,SJTG,SMVG,$JT,$MVG,SMVGO,PLB]{2,5}[0-9]{2,5}[A-Z]{0,1}\b", venue
    )
    # Fix venue errors during OCR runtime
    for i in term:
        i = i.replace("SMVGO", "SMVG")
        i = i.replace("$", "S")
        if i.startswith("STS"):
            venue = term[1]
        else:
            venue = i
    return venue


def fetch_data(image):
    """Function to fetch timetable from image"""
    # Gel all pixels in the image - where BGR = (51,255,204), OpenCV colors order is BGR not RGB (green color)
    gray = np.all(image == (51, 255, 204), 2)  # gray is a logical matrix with True

    # Convert logical matrix to uint8
    gray = gray.astype(np.uint8) * 255
    cnts = cv2.findContours(gray, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_NONE)[-2]
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
        image = crop
        text = pt.image_to_string(image, lang="eng")
        slot = re.findall(r"^[A-Z]{1,3}[0-9]{1,2}", text)[0]
        course_name = re.findall(r"[A-Z]{3}[0-9]{4}", text)[0]
        course_type = re.findall(r"[ETH,SS,ELA,LO]{2,3}\b", text)
        course_type = "Lab" if course_type[0] in ("ELA", "LO") else "Theory"
        try:
            venue = get_venue(text)
        except IndexError:
            venue = None
        # writing data to json
        slot_data = {
            "Slot": slot,
            "Course_Name": course_name,
            "Course_type": course_type,
            "Venue": venue,
        }
        data.append(slot_data)

    write_json(data)
    # returning data
    return {"Slots": data}
