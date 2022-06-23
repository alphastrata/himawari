from PIL import Image
from tqdm import tqdm
import cv2
import numpy as np
import os
import sys


def get_images_list_comp(path="tiles"):
    return sorted(list(os.listdir(path)), key=lambda c: c[-10:])


def check(im):
    return im.shape[0] == im.shape[1]


if __name__ == "__main__":
    # get yymmdd and hhmmss from sys args (if any provided)
    yymmdd = sys.argv[1] if len(sys.argv) > 1 else "20190608"
    hhmmss = sys.argv[2] if len(sys.argv) > 2 else "040000"

    completed = f"completed/{yymmdd}/"
    tiles_path = f"tiles/{yymmdd}/{hhmmss}/"

    idx = 0

    tiles_sorted = get_images_list_comp(tiles_path)
    print("tiles found:", len(tiles_sorted))

    r2b = lambda rgb: (cv2.imread(rgb))
    verticals = []

    for ridx in tqdm(range(20)):
        outvec = []
        for cidx in range(20):
            tilepath = f"{os.path.join(tiles_path)}{hhmmss}R{ridx}_C{cidx}.png"
            outvec.append(r2b(tilepath))

        verticals.append(cv2.hconcat(outvec))

    # write out our completed image
    vstack = cv2.vconcat(verticals)

    outname = tiles_path.split("/")[-1]
    resname = os.path.join(completed, f"full-disc {hhmmss}.jpg")

    cv2.imwrite(f"{resname}", vstack)

    exit(0)
