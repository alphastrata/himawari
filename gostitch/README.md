```
                         __   ^
|-| ! |\/| /\ \ /\ / /\  |_!  |  
! ! | |  !/  \ \  / /  \ |\   !
                         | \    -気象衛星STITCHER
himawari-stitcher

```
## About:
This is a go implementation of the parent directory's [stitcher.py](../../stitcher.py), which concatenates the tiles from the scraper.

### Requirements:
1. `golang, progressbar/v3, gocv` [gocv](https://github.com/hybridgroup/gocv
2. a collection of himawari tiles, prepared by himawari
4. *this* repo
<p align="right">(<a href="#top">back to top</a>)</p>
 
## Usage:
_assuming you've cloned the repo and got all the above installed properly_
```
go build gostitch
./gostitch YYMMDD HHMMSS
```
_You may need to modify the run permissions, `chmod +x` on unix based systems._

<p align="right">(<a href="#top">back to top</a>)</p>
 
### Changelog
- [x] Take in YYYYMMDD HHMMSS as args like himawari does.
- [x] Flag for himawari so it calls this by default
- [x] Take in the tile path and desired path for the full discs to go as args.
- [x] Try out the loop with some go-rountines, see how many you can add before it blows up your pc...

<p align="right">(<a href="#top">back to top</a>)</p>


## Acknowledgments
* [CY](https://github.com/Subzerofusion)
* [hybridgroup](https://hybridgroup.com/) [twitter](https://twitter.com/GoCVio) [GH](https://github.com/hybridgroup/gocv)
* [Sean Doran](https://www.youtube.com/c/Se%C3%A1nDoran/videos)
* [NICT](https://www.nict.go.jp/index.html)

<p align="right">(<a href="#top">back to top</a>)</p>
