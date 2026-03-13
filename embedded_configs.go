package main

var embeddedConfigs = map[string]string{
	".bash_profile": `# Add ~/bin to PATH
export PATH="$HOME/bin:$PATH"

# Load modular shell config files
for file in ~/.{path,bash_prompt,exports,aliases,functions,extra}; do
    [ -r "$file" ] && [ -f "$file" ] && source "$file"
done
unset file

shopt -s nocaseglob
shopt -s histappend
shopt -s cdspell

for option in autocd globstar; do
    shopt -s "$option" 2>/dev/null
done

if which brew &>/dev/null && [ -r "$(brew --prefix)/etc/profile.d/bash_completion.sh" ]; then
    source "$(brew --prefix)/etc/profile.d/bash_completion.sh"
elif [ -f /etc/bash_completion ]; then
    source /etc/bash_completion
fi

if type _git &>/dev/null; then
    complete -o default -o nospace -F _git g
fi
`,

	".bash_prompt": `#!/usr/bin/env bash

prompt_git() {
    local s=''
    local branchName=''

    git rev-parse --is-inside-work-tree &>/dev/null || return

    branchName="$(git symbolic-ref --quiet --short HEAD 2>/dev/null || \
        git describe --all --exact-match HEAD 2>/dev/null || \
        git rev-parse --short HEAD 2>/dev/null || \
        echo '(unknown)')"

    if ! $(git diff --quiet --ignore-submodules --cached); then s+='+'; fi
    if ! $(git diff-files --quiet --ignore-submodules --); then s+='!'; fi
    if [ -n "$(git ls-files --others --exclude-standard)" ]; then s+='?'; fi
    if $(git rev-parse --verify refs/stash &>/dev/null); then s+='$'; fi

    [ -n "${s}" ] && s=" [${s}]"
    echo -e "${1}${branchName}${2}${s}"
}

if tput setaf 1 &>/dev/null; then
    tput sgr0
    bold=$(tput bold)
    reset=$(tput sgr0)
    white=$(tput setaf 15)
    yellow=$(tput setaf 136)
    green=$(tput setaf 64)
    violet=$(tput setaf 61)
else
    bold=''
    reset="\e[0m"
    white="\e[1;37m"
    yellow="\e[1;33m"
    green="\e[1;32m"
    violet="\e[1;35m"
fi

PS1="\[\033]0;\W\007\]"
PS1+="\[${bold}\]\n"
PS1+="\[${yellow}\]\u"
PS1+="\[${white}\] at \[${yellow}\]\h"
PS1+="\[${white}\] in \[${green}\]\w"
PS1+="\$(prompt_git \"\[${white}\] on \[${violet}\]\" \"\[${green}\]\")"
PS1+="\n\[${white}\]\\$ \[${reset}\]"
export PS1
`,

	".aliases": `#!/usr/bin/env bash

alias gst="git status"
alias gaa="git add --all"
alias ga="git add "
alias gc="git commit "
alias gco="git checkout "
alias gb="git branch"
alias gbv="git branch -vv"

alias cls="clear"
alias grep="grep --color=auto"
alias mv='mv -v'
alias mkdir='mkdir -v'
alias cp='cp -v'
alias rm='rm -v'
alias ln='ln -v'
alias ls='ls --color=auto'
alias ll='ls -Grlth'
alias ..='cd ../'
alias ...='cd ../../'

if command -v vim >/dev/null 2>&1; then
    alias vi="vim --noplugin"
    alias cvi="vim ~/.vim/vimrc"
fi
`,

	".exports": `#!/usr/bin/env bash

export EDITOR='vim'
export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8
export WORKON_HOME='~/.virtualenvs'
export PATH="$WORKON_HOME/bin:$PATH"
export XDG_CONFIG_HOME="${HOME}/.config"
export HISTSIZE='32768'
export GPG_TTY=$(tty)
export HISTIGNORE="history:ls:pwd:"
export HISTTIMEFORMAT='%F, %T '
export TERM='xterm-256color'
export TMUX_TMPDIR='~/.tmux/tmp'
`,

	".functions": `#!/usr/bin/env bash

r() { grep "$1" ${@:2} -n -R .; }
src() { source ~/.bashrc; }

with_proxy() {
    HTTPS_PROXY=socks5://127.0.0.1:1080 HTTP_PROXY=socks5://127.0.0.1:1080 "$@"
}
`,

	".gitconfig": `# ~/.gitconfig
[core]
    editor = vim
[user]
    email = 19486826+killua525@users.noreply.github.com
    name = killua525
[push]
    default = simple
[grep]
    extendRegexp = true
    lineNumber = true
[color]
    ui = auto
[commit]
    verbose = true
    gpgsign = true
[help]
    autocorrect = 1
[alias]
    aa = add -A
    ca = commit -vam
    ci = commit -v
    st = status
    glg = log --oneline --decorate --all --graph
    amend = commit --amend --reuse-message=HEAD
[difftool]
    prompt = false
[diff]
    tool = vimdiff
[merge]
    ff = only
[tig]
    vertical-split = true
    tab-size = 4
[tag]
    forcesignannotated = true
`,

	".gitignore": `# Compiled Python files
*.pyc

# Folder view configuration files
.DS_Store
Desktop.ini

# Thumbnail cache files
._*
Thumbs.db

# Files that might appear on external disks
.Spotlight-V100
.Trashes
`,

	".tigrc": `set tab-size = 4
set vertical-split = true

set main-view = date author:abbreviated id:yes,color \\
                line-number:no,interval=1 \\
                commit-title:graph=v2,refs=yes,overflow=no
set tree-view = date:default author:abbreviated id:yes,color \\
                line-number:no,interval=5 \\
                mode file-size:units,width=0 \\
                file-name

set blame-view-id = yes,color
set blame-view-line-number = yes,interval=1
set main-view-date = custom
set main-view-date-format = "%Y-%m-%d %H:%M"
set ignore-space = all
set log-options = --show-signature
set diff-options = --show-signature
`,

	".tmux.conf": `run-shell 'tmux setenv -g TMUX_VERSION $(tmux -V | sed -En "s/^tmux[^0-9]*([.0-9]+).*/\\1/p")'
set -sg escape-time 0
set -g prefix C-a
set -g default-terminal "screen-256color"

set -g base-index 1
setw -g pane-base-index 1
set -g renumber-windows on
set -g allow-rename off
set -g history-limit 5000

unbind %
bind | split-window -h
bind _ split-window -v
bind C-a send-prefix
bind a last-window

is_vim='echo "#{pane_current_command}" | grep -iqE "(^|/)(g?(view|n?vim?)(diff)?|git)$"'
bind -n C-h if-shell "$is_vim" "send-keys C-h" "select-pane -L"
bind -n C-j if-shell "$is_vim" "send-keys C-j" "select-pane -D"
bind -n C-k if-shell "$is_vim" "send-keys C-k" "select-pane -U"
bind -n C-l if-shell "$is_vim" "send-keys C-l" "select-pane -R"

bind r source-file ~/.tmux.conf \; display "Reloaded ~/.tmux.conf"
if-shell 'test -z "$SSH_CLIENT"' "source-file ~/.tmux-theme.conf"
set -g display-panes-time 1500
`,

	".tmux-theme.conf": `if-shell -b '[ "$(echo "$TMUX_VERSION > 2.1" | bc)" = 1 ]' {
    set -g status-bg colour234
    set -g status-fg white
    setw -g window-status-current-style bg=colour234
    setw -g window-status-current-style fg=yellow

    set -g pane-border-style fg=colour237
    set -g pane-active-border-style fg=colour221

    set -g status-interval 2
    set -g status-left-length 55
    set -g status-right-length 150
    set -g status-left '#[fg=blue]#S #[fg=brightblack]•#[default]'
    set -g status-right '#22T #[fg=blue]%H:%M#[default]'
}
`,

	".vim/vimrc": `" Main Vim configuration
set nocompatible
let mapleader=','
set cursorline
set title
set autoread
au FocusGained,BufEnter * checktime

set wildignore+=*.o,*~,*.pyc,*.class,*.swp,*.bak
set wildignore+=*/.git/*,*/.hg/*,*/.svn/*

set backspace=indent,eol,start
set whichwrap+=<,>,h,l
set smartcase
set hlsearch
set incsearch

set lazyredraw
set showmatch
set noerrorbells
set novisualbell
set t_vb=

syntax enable
set background=dark
colorscheme peaksea
set colorcolumn=81

set encoding=utf-8
set fileencodings=utf-8,gbk,gb18030,latin1
set termencoding=utf-8
set ffs=unix,dos,mac

set nobackup
set nowb
set noswapfile
set hidden

set expandtab
set smarttab
set shiftwidth=4
set tabstop=4
set autoindent
set smartindent

set foldcolumn=1
set cmdheight=2
set tags=./.tags;,.tags
set t_ti= t_te=
set completeopt-=preview
filetype indent on

autocmd BufReadPost quickfix nnoremap <buffer> <CR> <CR>
autocmd FileType python setlocal et sta sw=4 sts=4

set statusline=%1*\%<%.50F\
set statusline+=%=%2*\%y%m%r%h%w\ %*
set statusline+=%3*\%{&ff}\[%{&fenc}]\ %*
set statusline+=%4*\ row:%l/%L,col:%c\ %*
set statusline+=%5*\%3p%%\%*
set laststatus=2

aug QFClose
  au!
  au WinEnter * if winnr('$') == 1 && &buftype == "quickfix"|q|endif
aug END

try
    source <sfile>:p:h/functions.vim
    source <sfile>:p:h/keys.vim
    source <sfile>:p:h/plugin.vim
catch
endtry
`,

	".vim/plugin.vim": `call plug#begin('~/.vim/plugged')
Plug 'SirVer/ultisnips', { 'for' : ['cpp','go','markdown'] }
Plug 'honza/vim-snippets', { 'for' : ['cpp','go','markdown'] }
Plug 'Shougo/echodoc.vim', { 'for' : ['cpp','go'] }
Plug 'preservim/nerdcommenter'
Plug 'preservim/nerdtree'
Plug 'Yggdroot/LeaderF'
Plug 'fatih/vim-go', { 'for' : ['go'], 'do': ':GoUpdateBinaries' }
Plug 'mileszs/ack.vim', { 'for' : ['cpp','go'] }
Plug 'rhysd/vim-clang-format', { 'for' : ['cpp'] }
Plug 'editorconfig/editorconfig-vim'
Plug 'godlygeek/tabular'
Plug 'plasticboy/vim-markdown'
call plug#end()

let g:Lf_UseCache=0
let g:Lf_UseVersionControlTool=0
let g:Lf_ShowDevIcons = 0
let g:Lf_ShortcutF = '<C-P>'
let g:Lf_WindowPosition = 'popup'
let g:Lf_WindowHeight=0.3
nnoremap <leader>c :LeaderfCommand<CR>

augroup go
  au!
  let g:go_referrers_mode = 'gopls'
  let g:go_implements_mode = 'gopls'
  let g:go_rename_command = 'gopls'
  let g:go_fmt_command = 'gopls'
  let g:go_diagnostics_enabled = 1
  let g:go_imports_autosave = 1
augroup END

let g:NERDTrimTrailingWhitespace = 1
let g:NERDSpaceDelims = 1
let g:NERDDefaultAlign = 'left'
nnoremap <C-t> :NERDTreeToggle<CR>
let NERDTreeQuitOnOpen=1
`,

	".vim/gvimrc": `if has("gui_running")
    if has("mac") || has("macunix")
        set gfn=IBM\ Plex\ Mono\ 14,:Hack\ 14,Source\ Code\ Pro\ 12
    elseif has("linux")
        set gfn=IBM\ Plex\ Mono\ 14,:Hack\ 14,Source\ Code\ Pro\ 12
    elseif has("unix")
        set gfn=Monospace\ 11
    endif
    set guioptions-=T
    set guioptions-=e
    set t_Co=256
endif
`,

	".vim/.ycm_extra_conf.py": `import os
import os.path
import ycm_core

BASE_FLAGS = [
    '-Wall',
    '-Wextra',
    '-Werror',
    '-fexceptions',
    '-std=c++17',
    '-xc++',
    '-I/usr/lib/',
    '-I/usr/include/',
]

SOURCE_EXTENSIONS = ['.cpp', '.cxx', '.cc', '.c', '.m', '.mm']
HEADER_EXTENSIONS = ['.h', '.hxx', '.hpp', '.hh']


def IsHeaderFile(filename):
    extension = os.path.splitext(filename)[1]
    return extension in HEADER_EXTENSIONS


def MakeRelativePathsInFlagsAbsolute(flags, working_directory):
    if not working_directory:
        return list(flags)
    new_flags = []
    make_next_absolute = False
    path_flags = ['-isystem', '-I', '-iquote', '--sysroot=']

    for flag in flags:
        new_flag = flag
        if make_next_absolute:
            make_next_absolute = False
            if not flag.startswith('/'):
                new_flag = os.path.join(working_directory, flag)

        for path_flag in path_flags:
            if flag == path_flag:
                make_next_absolute = True
                break

            if flag.startswith(path_flag):
                path = flag[len(path_flag):]
                new_flag = path_flag + os.path.join(working_directory, path)
                break

        if new_flag:
            new_flags.append(new_flag)
    return new_flags


def FlagsForFile(filename):
    return {
        'flags': BASE_FLAGS,
        'do_cache': True,
    }
`,

	".vim/colors/peaksea.vim": `" Peaksea colorscheme placeholder generated by cfg_sync
set background=dark
hi clear
if exists("syntax_on")
  syntax reset
endif
let g:colors_name = "peaksea"
`,

	"config/sshrc.d/vimrc": `let mapleader="," 
set path+=**
set wildignore+=*.o,*~,*.pyc,*.class,*.swp
set nocompatible
set wildmenu
set noswapfile
set ts=2
set expandtab
set shiftwidth=2
set softtabstop=2
set incsearch
set hlsearch
set showmatch
set backspace=2
set cursorline
set encoding=utf-8
set textwidth=100
set colorcolumn=101
set cmdheight=2
set updatetime=300
set background=light
colo desert

cnoremap <C-a> <Home>
cnoremap <C-e> <End>
noremap H ^
noremap L $
nnoremap <leader>q :q<CR>
nnoremap <leader>w :w<CR>

set hidden
set t_ti= t_te=

nnoremap <F2> :set nu! rnu!<cr>
nnoremap <F3> :set list!<cr>
nnoremap <F5> :set paste!<cr>
`,

	"config/sshrc.d/.vimrc": `let mapleader="," 
set path+=**
set wildignore+=*.o,*~,*.pyc,*.class,*.swp
set nocompatible
set wildmenu
set noswapfile
set ts=2
set expandtab
set shiftwidth=2
set softtabstop=2
set incsearch
set hlsearch
set showmatch
set backspace=2
set cursorline
set encoding=utf-8
set textwidth=100
set colorcolumn=101
set cmdheight=2
set updatetime=300
`,
}
