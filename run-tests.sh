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

Today=$(date +'%d-%m-%Y')
Yesterday=$(date -d "-1 days" +"%d-%m-%Y")


# Simple
echo -e "\n${Cyan}go test -v ${Color_Off}\n"
go test -v

# Run benchmark tests
NewFname="$Today.bench"
echo -e "\n${Cyan}go test -v -bench=. -benchmem ${Color_Off} > ${NewFname}\n"
go test -bench=. -benchmem > $NewFname

# Run benchmark tests with data race check
echo -e "\n${Cyan}go test -v -race -bench=. -benchmem ${Color_Off}\n"
go test -race -bench=. -benchmem

# Compare to bench profiles
OldFname="$Yesterday.bench"
echo -e "\n\n${Yellow}"
benchcmp $OldFname $NewFname
echo -e "${Color_Off}\n\n"
