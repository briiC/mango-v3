# Reset color
Color_Off='\e[0m'       # Text Reset

# Regular Colors
Black='\e[0;30m'        # Black
Red='\e[0;31m'          # Red
Green='\e[0;32m'        # Green
Yellow='\e[0;33m'       # Yellow
Blue='\e[0;34m'         # Blue
Purple='\e[0;35m'       # Purple
Cyan='\e[0;36m'         # Cyan
White='\e[0;37m'        # White

FNAME="cover.log"

# Collect data
echo -e "\n${Cyan} go test -coverprofile=$FNAME ${Color_Off}\n"
    go test -coverprofile=$FNAME


# Analyze & display data
echo -e "\n${Cyan} go tool cover -func=$FNAME ${Color_Off}\n"
    go tool cover -func=$FNAME

echo -e "\n### Analyze it by yourself using -html flag:"
echo -e "### ${Purple}go tool cover -html=$FNAME ${Color_Off}\n"
