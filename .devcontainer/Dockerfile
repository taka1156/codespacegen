FROM golang:1.26

WORKDIR /workspace

ENV LANG=ja_JP.UTF-8 \
    LANGUAGE=ja_JP:ja \
    LC_ALL=ja_JP.UTF-8 \
    TZ=Asia/Tokyo

RUN <<-EOF
curl -1sLf 'https://dl.cloudsmith.io/public/evilmartians/lefthook/setup.deb.sh' | bash
apt update
apt install -y locales tzdata git vim curl shfmt lefthook
apt clean
EOF

RUN <<-EOF
echo "ja_JP.UTF-8 UTF-8" > /etc/locale.gen
locale-gen ja_JP.UTF-8
ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
dpkg-reconfigure -f noninteractive tzdata
EOF

RUN curl -fsSL \
  https://github.com/shipengqi/commitizen/releases/download/v2.3.5/commitizen-v2.3.5-linux-amd64.tar.gz \
  | tar -xz -C /usr/local/bin commitizen \
  && chmod +x /usr/local/bin/commitizen \
  && commitizen init

RUN curl -fsSL \
  https://github.com/conventionalcommit/commitlint/releases/download/v0.12.0/commitlint_v0.12.0_linux_amd64.tar.gz \
  | tar -xz -C /usr/local/bin commitlint \
  && chmod +x /usr/local/bin/commitlint

RUN curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.11.4

RUN <<-EOF
curl -o ~/.git-prompt.sh https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh
curl -o ~/.git-completion.sh https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash

echo "[ -f ~/.git-prompt.sh ] && source ~/.git-prompt.sh" >> ~/.bashrc
echo "[ -f ~/.git-completion.sh ] && source ~/.git-completion.sh" >> ~/.bashrc
echo "GIT_PS1_SHOWDIRTYSTATE=true" >> ~/.bashrc
echo "GIT_PS1_SHOWUNTRACKEDFILES=true" >> ~/.bashrc
echo "GIT_PS1_SHOWUPSTREAM=auto"  >> ~/.bashrc
git config --system --add safe.directory /workspace
echo 'export PS1="\[\033[01;32m\]\u@\h\[\033[01;33m\] \w \[\033[01;31m\]\$(__git_ps1 \"(%s)\") \\n\[\033[01;34m\]\\$ \[\033[00m\]"' >> ~/.bashrc
EOF

CMD ["bash"]
