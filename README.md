```
                         __   ^
|-| ! |\/| /\ \ /\ / /\  |_!  |  
! ! | |  !/  \ \  / /  \ |\   !
                         | \    -気象衛星
himawari

```
## About:

[himawari](https://himawari8.nict.go.jp/) is a Japanese satellite for weather monitoring, it takes an image of the full earth's disc every ten minutes, and has done so since 2015.
The highest available resolution of the tiles the disc image made up of is 550 by 550px, so when you stitch all 400 images together to make a full disc you get a whopping 11000 by 11000px, a 121MP image. 

A complete day's worth of discs is about ~11.5GB..."https://imgur.com/a/TAjJ1mG", you can see the example `full_disc_example.png` from this repo.

## Requirements:
1. `python3, opencv-python, tqdm, numpy, PIL, golang, progressbar/v3, gocv` [gocv](https://github.com/hybridgroup/gocv
3. internet
4. *this* repo `git clone https://github.com/alphastrata/himawari`
5. enough of an understanding of how to work this sorta thing...
6. If you're using batch mode you'll need a commandline configurable vpn. (this example uses PIA)
7. patience
<p align="right">(<a href="#top">back to top</a>)</p>
 
# Usage:
_assuming you've cloned the repo and got all the above installed properly_
### Oneshots:
```
cd <wherever you put himawari>
go mod init himawari
go mod tidy 
```
To prepare the dependencies etc... then,
```
go run himawari 20180818 225000
```
or 
```
go build himawari
./himawari 20180818 225000
```
_This will take a single oneshot, worth of tiles and save it to `./tiles/YYYYMMDD/HHMMSS` of 22:50 on the 18th of August 2018._

### Batch:
```
cd <wherever you put himawari>
nvim runner.py # replace nvim with the text-editor of choice
```
In `runner.py` configure your date range, line 183 etc. Configure your VPN on `189` and `214`.
```
pip install -r requirements.txt 
python3 runner.py
```
If you're unsure about your system's python version or command -- what are you even doing here?

_This will take a range you specify, be patient for long ranges, on avg a full disc takes between 11 to 150s to download, stitch and save. There are 144 potential discs in a day._

### /examples
In this dir you'll find a bunch of helper scripts to help you shrink images, stitch tiles into complete discs or rename tilesets.

Now you just wait...
<p align="right">(<a href="#top">back to top</a>)</p>
 

## Acknowledgments
* [CY](https://github.com/Subzerofusion)
* [Sean Doran](https://www.youtube.com/c/Se%C3%A1nDoran/videos)
* [NICT](https://www.nict.go.jp/index.html)
<p align="right">(<a href="#top">back to top</a>)</p>
