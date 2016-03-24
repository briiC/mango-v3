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


# # Simple
# echo -e "\n${Cyan}go test -v ${Color_Off}\n"
# go test -v
#
# # Simple with data recing
# echo -e "\n${Cyan}go test -v -race ${Color_Off}\n"
# go test -v -race

# Run benchmark tests
echo -e "\n${Cyan}go test -v -bench=. ${Color_Off}\n"
go test -v -bench=.

# Run benchmark tests with data race check
echo -e "\n${Cyan}go test -v -race -bench=. ${Color_Off}\n"
go test -v -race -bench=.
