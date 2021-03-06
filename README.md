# HCCTV-server

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Mysql](https://img.shields.io/badge/MySQL-005C84?style=for-the-badge&logo=mysql&logoColor=white)
![Nginx](https://img.shields.io/badge/Nginx-009639?style=for-the-badge&logo=nginx&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Azure](https://img.shields.io/badge/Azure_DevOps-0078D7?style=for-the-badge&logo=azure-devops&logoColor=white)

| <a id="a1"></a>목차         |
| --------------------------- |
| [1. 프로젝트 init](#1)<br/> |
| [2. 브랜치 ](#2)<br/>       |

<br/>

# <a id="1"></a>[1](#a1). 프로젝트 init

> ---
>
> # docker
>
> - docker가 없다면 docker 설치 후 진행해주세요
>   - osx
>   ``` 
>   brew install docker 
>   ```
>   - linux
>   ```
>   sudo apt-get install docker-engine -y && sudo apt-get install docker-engine -y
>   ```
>
> # local 세팅
>
>       git clone https://github.com/WhiteHCCTV/HCCTV-server
>
> # 서버 구동
> ## run.sh deploy
>
> - 배포 서버를 위한 구동
> - Azure, domain setting으로 서버가 구동됩니다 ( 추후 구현 )
>
> ## run.sh dev
>
> - 로컬 개발 서버를 위한 구동
> - localhost db, domain setting으로 서버가 구동됩니다.
>
> ---

# <a id="2"></a>[2](#a1). 브랜치
> - main : release 
>   - hotfix
> - develop : feature / patch 
>   - feature
>   - patch
# <a id="3"></a>[3](#a1). 주의사항

# <a id="3"></a>[4](#a1). Trouble shooting
