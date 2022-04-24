"""Import necessary modules"""

import cv2
import numpy as np
import uvicorn
from fastapi import FastAPI, File, Form, HTTPException, UploadFile
from starlette.middleware.cors import CORSMiddleware
from tableDetection import fetch_text_timetable

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://127.0.0.1:8000",
        "http://0.0.0.0:8000",
        "http://vittyapi.dscvit.com",
        "https://vittyapi.dscvit.com",
        "https://vitty.pages.dev",
        "https://vitty.dscvit.com",
        "http://vitty.dscvit.com"
    ],  # Allows all origins
    allow_credentials=True,
    allow_methods=["GET", "POST"],
    allow_headers=["*"],
)



@app.get("/")
async def testing():
    """Check if server is working"""
    return "Ok! Working!"


@app.post("/uploadtext/")
async def get_timetable(request: str = Form(...)):
    data = fetch_text_timetable(request)
    return data


if __name__ == "__main__":
    uvicorn.run(app)
