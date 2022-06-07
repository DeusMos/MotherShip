FROM ubuntu:20.04
ENV DEBIAN_FRONTEND noninteractive
RUN useradd -rm -d /home/xp -s /bin/bash -g root -G sudo xp
RUN useradd -rm -d /home/user -s /bin/bash -g root -G sudo user
RUN echo 'root:HK$YhxPGxbe)@P7?' | chpasswd
RUN echo 'user:password' | chpasswd
RUN echo 'xp:password' | chpasswd
RUN apt-get update &&  apt-get install -y sudo openssh-server nano nmap
RUN mkdir /var/run/sshd
RUN sed -ri 's/^#?PermitRootLogin\s+.*/PermitRootLogin yes/' /etc/ssh/sshd_config
# RUN sed -ri 's/UsePAM yes/#UsePAM yes/g' /etc/ssh/sshd_config
RUN sed -ri 's/#GatewayPorts no/GatewayPorts yes/' /etc/ssh/sshd_config
RUN sed -ri 's/#Port 22/Port 2222/' /etc/ssh/sshd_config
# RUN echo 'alias connectg6="ssh -o StrictHostKeyChecking=no -o \"UserKnownHostsFile /dev/null\" -p 2222 user@localhost"' >> ~/.bashrc
# RUN echo 'alias connectxp="ssh -o StrictHostKeyChecking=no -o \"UserKnownHostsFile /dev/null\" -p 2222 xp@localhost"' >> ~/.bashrc
CMD service ssh start && tail -f /dev/null