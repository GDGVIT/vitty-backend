<p align="center">
<a href="https://dscvit.com">
	<img src="https://user-images.githubusercontent.com/30529572/92081025-fabe6f00-edb1-11ea-9169-4a8a61a5dd45.png" alt="DSC VIT"/>
</a>
	<h2 align="center"> Vitty Backend</h2>
	<h4 align="center"> An application that extracts and reads data from a timetable and sends push notifications to user to attend scheduled classes.</h4>
</p>

---
[![Join Us](https://img.shields.io/badge/Join%20Us-Developer%20Student%20Clubs-red)](https://dsc.community.dev/vellore-institute-of-technology/)
[![Discord Chat](https://img.shields.io/discord/760928671698649098.svg)](https://discord.gg/498KVdSKWR)

[![DOCS](https://img.shields.io/badge/Documentation-see%20docs-green?style=flat-square&logo=appveyor)](https://vittyapi.dscvit.com/docs) 
  [![UI ](https://img.shields.io/badge/User%20Interface-Link%20to%20UI-orange?style=flat-square&logo=appveyor)](https://www.figma.com/file/3ILW1qy1qIjiJ5S78zyIqh/VITTY?node-id=1%3A4)


## Features
- [X]  Extracts timetable from the image
- [X]  Manually copy and paste timetable
<!-- - [ ]  < feature >
- [ ]  < feature > -->

<br>

## Dependencies

- FastAPI
- Numpy
- OpenCV
- Pytesseract
- Uvicorn


## Running

### Directions to Install
```bash
pipenv install
pipenv shell
uvicorn main:app
```

### Directions to Execute

```bash
docker build .
docker run -p 80:8000 -d <IMAGE_NAME>
```
<br>

## Contributors

<table>
	<tr align="center">
		<td>
		Vishesh Bansal
		<p align="center">
			<img src = "https://avatars.githubusercontent.com/u/22132836" width="150" height="150" alt="Your Name Here (Insert Your Image Link In Src">
		</p>
			<p align="center">
				<a href = "https://github.com/VisheshBansal">
					<img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/>
				</a>
				<a href = "https://www.linkedin.com/in/bansalvishesh" target="_blank">
					<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
				</a>
			</p>
		</td>
	</tr>
</table>

<p align="center">
	Made with :heart: by <a href="https://dscvit.com">DSC VIT</a>
</p>
