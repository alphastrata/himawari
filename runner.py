import os
from random import randint
import time
import logging


TIMESTAMP_LIST = [
    "160000",
    "161000",
    "162000",
    "163000",
    "164000",
    "165000",
    "170000",
    "171000",
    "172000",
    "173000",
    "174000",
    "175000",
    "180000",
    "181000",
    "182000",
    "183000",
    "184000",
    "185000",
    "190000",
    "191000",
    "192000",
    "193000",
    "194000",
    "195000",
    "200000",
    "201000",
    "202000",
    "203000",
    "204000",
    "205000",
    "210000",
    "211000",
    "212000",
    "213000",
    "214000",
    "215000",
    "220000",
    "221000",
    "222000",
    "223000",
    "224000",
    "225000",
    "230000",
    "231000",
    "232000",
    "233000",
    "234000",
    "235000",
    "000000",
    "001000",
    "002000",
    "003000",
    "004000",
    "005000",
    "010000",
    "011000",
    "012000",
    "013000",
    "014000",
    "015000",
    "020000",
    "021000",
    "022000",
    "023000",
    "024000",
    "025000",
    "030000",
    "031000",
    "032000",
    "033000",
    "034000",
    "035000",
    "040000",
    "041000",
    "042000",
    "043000",
    "044000",
    "045000",
    "050000",
    "051000",
    "052000",
    "053000",
    "054000",
    "055000",
    "060000",
    "061000",
    "062000",
    "063000",
    "064000",
    "065000",
    "070000",
    "071000",
    "072000",
    "073000",
    "074000",
    "075000",
    "080000",
    "081000",
    "082000",
    "083000",
    "084000",
    "085000",
    "090000",
    "091000",
    "092000",
    "093000",
    "094000",
    "095000",
    "100000",
    "101000",
    "102000",
    "103000",
    "104000",
    "105000",
    "110000",
    "111000",
    "112000",
    "113000",
    "114000",
    "115000",
    "120000",
    "121000",
    "122000",
    "123000",
    "124000",
    "125000",
    "130000",
    "131000",
    "132000",
    "133000",
    "134000",
    "135000",
    "140000",
    "141000",
    "142000",
    "143000",
    "145000",
    "150000",
    "151000",
    "152000",
    "153000",
    "154000",
    "155000",
]


VPNS = [
    "au-perth",
    "au-melbourne",
    "au-sydney",
    "new-zealand",
    "taiwan",
    "hong-kong",
    "singapore",
    "jp-tokyo",
]

# TODO: generate the list like instead of that god-awful
# MAXMINS     = 60
# MAXHOURS    = 24
# MIN_INC     = 10

YEAR = 2022
MONTH = 1
START_DATE = 1
END_DATE = 21
hh = 0
mm = 0

if __name__ == "__main__":
    start_timer = time.time()
    logging.basicConfig(
        filename="runner.log",
        format="%(asctime)s %(message)s",
        datefmt="%d/%b/%y %H:%M:%S",
        level=logging.INFO,
    )
    os.system("piactl connect")
    logging.info("VPN ON!")
    logging.info(f"STARTTIME: {start_timer}")

    # Go over the list of timestamps, use an external call to himawari to download them.
    # himawari will handle the download and external calls to `stitcher.py` which stitches
    # the tiles together.
    while START_DATE <= END_DATE:
        for hms in TIMESTAMP_LIST:
            pyargs_complete = f"{YEAR}{MONTH:0>2}{START_DATE:0>2} {hms}"  # TODO: fix the MONTH and day padding
            os.system("go run himawari.go " + pyargs_complete)
            logging.info(f"STITCHER: {pyargs_complete}")

        # NOTE: you may need to disable / edit this line to accomidate your VPN setup.
     new_rand = randint(0, len(VPNS) - 1)
        os.system(f"piactl set region {VPNS[new_rand]}")
        logging.info(f"VPN ON : {VPNS[START_DATE]}")

        START_DATE += 1
        time.sleep(2)

    logging.info(f"STARTDATE: {START_DATE}")
    logging.info(f"ENDDATE  : {END_DATE}")
    logging.info(f"TILES    : {END_DATE - START_DATE * 144}")
    os.system("piactl disconnect")
