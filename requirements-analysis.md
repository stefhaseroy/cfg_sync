# Requirement Analysis: killua525/dotfiles

## 1. Document Purpose
This document captures requirement analysis for the repository https://github.com/killua525/dotfiles.
The objective is to identify:
- Which software/tools are configured
- What configuration behaviors are defined for each tool
- Cross-platform installation and delivery expectations

## 2. Project Positioning
This repository is a cross-platform personal environment setup, targeting:
- macOS
- Linux (including WSL2)
- Windows (mainly Vim/gVim + Git workflow)

Core value:
- One repository provides unified shell, editor, terminal multiplexer, Git, and Kubernetes prompt experience.

## 3. Scope
### 3.1 In Scope
- Shell environment configuration (Bash)
- Vim/gVim configuration and plugin ecosystem
- Git and Tig configuration
- tmux behavior and theme configuration
- Kubernetes prompt (kube-ps1) module
- SSH remote Vim minimal profile (sshrc)
- Cross-platform bootstrap/update/install scripts

### 3.2 Out of Scope
- Business application code
- CI/CD pipeline engineering
- Deployment for production services

## 4. Software Configuration Inventory

### 4.1 Bash
Configuration files:
- .bash_profile
- .bash_prompt
- .aliases
- .exports
- .functions

Configured content:
- Startup loading chain for modular dotfiles
- Shell option tuning (history append, typo correction, recursive glob)
- Command completion enhancement (Git alias, SSH host, defaults, killall)
- Prompt with Git branch and workspace dirty-state indicators
- Large alias set for Git and common CLI commands
- Environment variable set (editor, locale, history, xdg, tmux temp path, gpg tty)
- Custom utility functions (quick bootstrap, proxy command wrapper, sync helper)

### 4.2 Vim
Configuration files:
- .vim/vimrc
- .vim/plugin.vim
- .vim/.ycm_extra_conf.py
- .vim/gvimrc
- .vim/colors/peaksea.vim

Configured content:
- Core editing behavior: encoding, search, indentation, statusline, quickfix behavior
- UI/visual behavior: colorscheme, colorcolumn, cursorline, bell handling
- Plugin management via vim-plug
- Plugin set for snippets, commenting, tree/file finder, Go/C++/Markdown workflows
- Go development defaults using gopls for diagnostics, rename, format, references
- C/C++ completion flags logic through YCM config, supporting compile_commands.json and include scanning
- gVim GUI-specific defaults (font, GUI options, initial dimensions)

### 4.3 Windows Vim/gVim Entry
Configuration files:
- tools/windows/_vimrc
- tools/windows/_gvimrc

Configured content:
- Resolve runtime from %USERPROFILE%/vimfiles
- Source common vimrc first, then GUI-specific settings
- Apply mouse, font, and window defaults for gVim on Windows

### 4.4 Git
Configuration files:
- .gitconfig
- .gitignore

Configured content:
- User identity, editor, push mode, grep and color output tuning
- Signed commit/tag oriented setup
- diff/merge tool behavior (vimdiff, ff policy)
- Extensive aliases for commit/log/PR-MR checkout/branch cleanup workflows
- Ignore OS-generated files and cross-platform metadata noise

### 4.5 Tig
Configuration file:
- .tigrc

Configured content:
- Tab width and vertical split defaults
- Main/tree/blame view layout customization with commit metadata visibility
- Date display customization
- Whitespace and signature display options in log/diff views

### 4.6 tmux
Configuration files:
- .tmux.conf
- .tmux-theme.conf

Configured content:
- Prefix change to Ctrl+a and low escape-time for Vim-friendly interaction
- Pane/window numbering and history behavior
- Pane split, copy/paste, and clipboard integration bindings
- Smart pane navigation with Vim-process awareness
- Theme definitions for status bar, pane border, activity, and mode colors
- Conditional theme loading when not in nested/remote SSH sessions

### 4.7 Kubernetes Prompt (kube-ps1)
Configuration file:
- config/kube-ps1.sh

Configured content:
- Prompt segment showing current Kubernetes context and namespace
- Symbol/prefix/suffix/separator/color customization
- Bash/Zsh shell-aware rendering and escape handling
- KUBECONFIG change detection and cache refresh
- kubeon/kubeoff controls (session and global toggles)

### 4.8 SSH Remote Vim Profile
Configuration files:
- config/sshrc.d/vimrc
- config/sshrc.d/.vimrc

Configured content:
- Lightweight Vim setup for remote sessions
- Minimal key mappings for navigation, save/quit, toggle options
- No-plugin friendly experience with practical editing defaults

## 5. Platform Installation and Delivery Requirements

### 5.1 macOS/Linux/WSL
Expected behavior:
- Link core dotfiles into HOME
- Link config children into ~/.config
- Provide update script for safe git fetch/pull workflow

Key scripts:
- tools/bootstrap.sh
- tools/update.sh
- install.sh

### 5.2 Windows
Expected behavior:
- Clone/update repo via PowerShell installer
- Link or copy config files into %USERPROFILE%
- Support fallback when symlink/junction is unavailable

Key scripts:
- tools/windows/setup.ps1
- tools/windows/_vimrc
- tools/windows/_gvimrc

## 6. Functional Requirement Summary
- FR-01: The solution shall provide modular shell configuration with prompt, aliases, exports, and utility functions.
- FR-02: The solution shall provide a complete Vim runtime including plugins and language-focused workflows.
- FR-03: The solution shall provide Git/Tig workflow acceleration through aliases and display customization.
- FR-04: The solution shall provide tmux interaction optimized for keyboard-heavy and Vim-centric navigation.
- FR-05: The solution shall provide Kubernetes context/namespace visibility in prompt with switchable state.
- FR-06: The solution shall support cross-platform installation and update for Unix-like and Windows environments.
- FR-07: The solution shall include remote-friendly lightweight Vim profile for SSH scenarios.

## 7. Non-Functional Requirements
- NFR-01: Configuration should be idempotent or safely repeatable by scripts.
- NFR-02: Existing user files should be backed up before replacement or relinking.
- NFR-03: Setup should degrade gracefully when optional dependencies are absent.
- NFR-04: UX consistency should be maintained across macOS/Linux/WSL/Windows as much as possible.

## 8. Acceptance Criteria
- AC-01: After installation, shell startup loads expected modules without errors.
- AC-02: Vim launches with base settings and plugin configuration available.
- AC-03: Git aliases and signing-related options are effective.
- AC-04: tmux key bindings and theme behavior match configured expectations.
- AC-05: kube-ps1 prompt segment can display context/namespace and be toggled on/off.
- AC-06: Windows installer can provision Vim/gVim and Git-related files into user profile.

## 9. Risks and Assumptions
Risks:
- Dependency availability differs by platform (for example xclip, gopls, kubectl).
- Terminal capability differences may affect color/font behavior.
- Symlink permissions on Windows may require elevated policy or fallback path.

Assumptions:
- Users have basic shell/Git environment available.
- Users accept personal-style key bindings and aliases.
- Optional tooling (kubectl, gopls, htop, prettyping, fasd, etc.) is installed when corresponding features are expected.

## 10. Suggested Next Step
Create a validation checklist script per platform to automatically verify each acceptance criterion after bootstrap.