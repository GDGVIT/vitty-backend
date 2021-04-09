'''Import necessary modules'''
import cv2
from fastapi import FastAPI, File, UploadFile
import numpy as np
import uvicorn
from tableDetection import detect_table

app = FastAPI()

@app.get("/test")
async def testing():
    """Check if server is working"""
    return "Ok Working"


@app.post("/predict/image")
async def predict_api(file: UploadFile = File(...)):
    """Upload timetable"""
    extension = file.filename.split(".")[-1] in ("jpg", "jpeg", "png", "JPG")
    if not extension:
        return {"Error": "Image must be jpg or png format!"}

    image = await file.read()
    nparr = np.fromstring(image, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    table = detect_table(img)
    return table


if __name__ == "__main__":
    uvicorn.run(app, debug=True)
