import json
import os

import cv2

from utils.api import ocr_space_file


# import imutils

# This only works if there's only one table on a page
# Important parameters:
#  - morph_size
#  - min_text_height_limit
#  - max_text_height_limit
#  - cell_threshold
#  - min_columns
def otsu(image):
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    thresh = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY_INV + cv2.THRESH_OTSU)[1]
    return thresh


def draw_contours(file, morph_params, thresh):
    param = cv2.getStructuringElement(cv2.MORPH_RECT, morph_params)
    detect_param = cv2.morphologyEx(thresh, cv2.MORPH_OPEN, param, iterations=2)
    cnts = cv2.findContours(detect_param, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    cnts = cnts[0] if len(cnts) == 2 else cnts[1]
    for c in cnts:
        cv2.drawContours(file, [c], -1, (36, 255, 12), 2)
    return file


def pre_process_image(img, save_in_file):
    # get rid of the color
    pre = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    yellow_hi = (70, 255, 255)
    yellow_lo = (40, 0, 0)
    # hsv = cv2.cvtColor(pre, cv2.COLOR_BGR2HSV)
    # Mask image to only select yellow
    mask = cv2.inRange(img, yellow_lo, yellow_hi)
    inv_mask = cv2.bitwise_not(mask)

    # Change image to black where we found yellow
    img[inv_mask > 0] = (0, 0, 0)
    thresh = otsu(img)
    # Detect horizontal lines
    draw_contours(img, (40, 1), thresh)
    # detect vertical lines
    draw_contours(img, (1, 10), thresh)

    if save_in_file is not None:
        cv2.imwrite(save_in_file, pre)
    return pre


def table_detection(image):
    # Load image, convert to grayscale, Otsu's threshold
    thresh = otsu(image)

    # Detect horizontal lines
    draw_contours(image, (40, 1), thresh)
    # detect vertical lines
    draw_contours(image, (1, 10), thresh)
    payload_image = os.path.join("data", "payload.png")
    cv2.imwrite(payload_image, image)
    return payload_image


def detect_table(image):
    pre_file = os.path.join("data", "pre.png")
    out_file = os.path.join("data", "out.png")
    img = image
    out_img = cv2.imread(os.path.join(out_file), cv2.IMREAD_COLOR)

    pre_process_image(img, pre_file)

    table_detection(out_img)
    # Visualize the result
    vis = img.copy()
    cv2.imwrite(out_file, vis)

    final_messages = []

    test_file = ocr_space_file(
        filename=out_file, overlay=True, language="eng"
    )

    test_file = json.loads(test_file)
    message = test_file["ParsedResults"][0]["ParsedText"]

    final_messages.append(message)

    print(final_messages)

    return {"Slots": final_messages}
