FROM ubuntu:24.04

WORKDIR /workspace

ENV LANG=en.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en.UTF-8 \
  TZ=Asia/Tokyo

RUN <<-EOF
apt update
apt install -y --no-install-recommends \
  bash \
  bash-completion \
  ca-certificates \
  tzdata \
  git \
  git-lfs \
  vim \
  curl \
  locales
rm -rf /var/lib/apt/lists/*
locale-gen en.UTF-8
update-locale LANG=en.UTF-8 LC_ALL=en.UTF-8
EOF

RUN git lfs install

RUN <<-EOF
ln -fs /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
dpkg-reconfigure -f noninteractive tzdata
EOF

RUN <<-EOF
curl -o ~/.git-prompt.sh https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh
curl -o ~/.git-completion.sh https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash

echo "[ -f ~/.git-prompt.sh ] && source ~/.git-prompt.sh" >> ~/.bashrc
echo "[ -f ~/.git-completion.sh ] && source ~/.git-completion.sh" >> ~/.bashrc
echo "GIT_PS1_SHOWDIRTYSTATE=true" >> ~/.bashrc
echo "GIT_PS1_SHOWUNTRACKEDFILES=true" >> ~/.bashrc
echo "GIT_PS1_SHOWUPSTREAM=auto"  >> ~/.bashrc
git config --system --add safe.directory /workspace
echo 'export PS1="\[\033[01;32m\]\u@\h\[\033[01;33m\] \w \[\033[01;31m\]\$(__git_ps1 \"(%s)\") \\n+\[\033[01;34m\]\\$ \[\033[00m\]"' >> ~/.bashrc
EOF

RUN curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash

CMD ["bash"]
