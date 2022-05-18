# =============================================================================
#                                TILE RENAMER
# =============================================================================

# Some older versions of himwari-scraper had a different naming convention, this script cleans them up.
# Old format:
#     /235000C16R9.png
#     /235000C16R10.png
# New format:
#     /160000R9_C13.png
#     /160000R8_C4.png

import os
import glob
from pathlib import Path
import logging

logging.basicConfig(
    filename="files_renamed.log",
    level=logging.INFO,
    format="%(asctime)s %(message)s",
    datefmt="%d/%b/%y %H:%M:%S",
)


def process_tiles():
    for root, dirs, files in os.walk("tiles"):
        for file in files:
            if file.endswith(".png"):
                tile_path = os.path.join(root, file)
                if is_tile_old(tile_path):

                    # adjust the name to the new format
                    new_name = tile_path.replace("000C", "a")
                    new_name = new_name.replace("R", "_C")
                    new_name = new_name.replace("a", "000R")

                    # get the absolute path ot the tile
                    os.rename(tile_path, new_name)

                    # log that we renamed the tile to files_renamed.log
                    logging.info(f"Renamed {tile_path} to {new_name}")


def is_tile_old(p) -> bool:
    # check if the tile is in the old format
    if "000C" in p:
        return True
    else:
        return False


if __name__ == "__main__":
    process_tiles()
    print("Done!")
