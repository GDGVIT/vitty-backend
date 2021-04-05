import cv2
import numpy as np


def convert_img():
    data = cv2.imread("../data/sample.PNG")
    info = np.iinfo(data.dtype)  # Get the information of the incoming image type
    data = data.astype(np.float64) / info.max  # normalize the data to 0 - 1
    print(info)
    data = 255 * data  # Now scale by 255
    img = data.astype(np.uint8)
    cv2.imwrite("../data/input.png", img)
