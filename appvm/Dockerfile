FROM ubuntu:20.04

# Install dependencies
RUN apt update \
    && apt-get update -y \
    && apt-get install --no-install-recommends --assume-yes wget software-properties-common gpg-agent supervisor xvfb mingw-w64 ffmpeg cabextract aptitude vim pulseaudio \
    && apt-get clean \
    && apt-get autoremove

# Install wine
ARG WINE_BRANCH="stable"
RUN wget -nv -O- https://dl.winehq.org/wine-builds/winehq.key | APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=1 apt-key add - \
    && apt-add-repository "deb https://dl.winehq.org/wine-builds/ubuntu/ $(grep VERSION_CODENAME= /etc/os-release | cut -d= -f2) main" \
    && dpkg --add-architecture i386 \
    && apt-get update \
    && DEBIAN_FRONTEND="noninteractive" apt-get install -y --install-recommends winehq-${WINE_BRANCH} \
    && rm -rf /var/lib/apt/lists/*

# Install winetricks
RUN wget -nv -O /usr/bin/winetricks https://raw.githubusercontent.com/Winetricks/winetricks/master/src/winetricks \
    && chmod +x /usr/bin/winetricks

# Download gecko and mono installers
COPY download_gecko_and_mono.sh /root/download_gecko_and_mono.sh
RUN chmod +x /root/download_gecko_and_mono.sh \
    && /root/download_gecko_and_mono.sh "$(dpkg -s wine-${WINE_BRANCH} | grep "^Version:\s" | awk '{print $2}' | sed -E 's/~.*$//')"

# Install graphic libraries
RUN winetricks d3dx9_43
# RUN winetricks --force -q dotnet48

# Silence all fixme warnings from wine
ENV WINEDEBUG fixme-all

WORKDIR /appvm

COPY default.pa /etc/pulse/
COPY supervisord.conf /etc/supervisor/conf.d/
COPY syncinput.cpp ./syncinput.cpp

# Compile syncinput.cpp
RUN x86_64-w64-mingw32-g++ ./syncinput.cpp -o ./syncinput.exe -lws2_32 -lpthread -static

# Add reg files to wineprefix
ARG APP_NAME
COPY apps/${APP_NAME}/setup/ /root/.wine/

ENTRYPOINT ["supervisord"]