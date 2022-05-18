from PIL import Image
from tqdm import tqdm
import cv2
import numpy as np
import os
import sys


def get_images_list_comp(path="tiles"):
    return sorted(list(os.listdir(path)), key=lambda c: c[-10:])


def build_from_stack(stack, filename):
    return cv2.hconcat(stack)


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

    r2b = lambda rgb: (cv2.imread(rgb))
    verticals = []

    for idx in tqdm(range(20)):
        outvec = []
        row_idx = f"R{idx}.png"

        for i in range(len(tiles_sorted)):
            if row_idx in tiles_sorted[i]:
                img = r2b(os.path.join(tiles_path, tiles_sorted[i]))
                outvec.append(img)

        filename = f"{row_idx}"

        verticals.append(build_from_stack(outvec, filename))

    # write out our completed image
    vstack = cv2.vconcat(verticals)
    cv2.imwrite(os.path.join(completed, f"{hhmmss}.png"), vstack)
