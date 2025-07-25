echo "This process will install Golang, G++, Rust, and any compilers/interpreters required by MK-TAN_INFL.

This process assumes the APT package manager is installed & functional, and you have sudo privileges. To cancel, click CTRL+C within 15 seconds."
sleep 15
sudo apt install g++ golang rustc python3
echo "Done!"
