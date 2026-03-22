#!/bin/sh
# check-path.sh — 설치 경로가 PATH에 없으면 안내
INSTALL_DIR="$1"

case ":$PATH:" in
    *":$INSTALL_DIR:"*)
        # PATH에 있음 — 안내 불필요
        ;;
    *)
        SHELL_RC=""
        if [ -n "$ZSH_VERSION" ] || [ "$(basename "$SHELL")" = "zsh" ]; then
            SHELL_RC="~/.zshrc"
        elif [ -f "$HOME/.bashrc" ]; then
            SHELL_RC="~/.bashrc"
        elif [ -f "$HOME/.bash_profile" ]; then
            SHELL_RC="~/.bash_profile"
        else
            SHELL_RC="~/.profile"
        fi

        echo ""
        echo "⚠ $INSTALL_DIR is not in your PATH."
        echo "  Add it with:"
        echo ""
        echo "    echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> $SHELL_RC"
        echo "    source $SHELL_RC"
        echo ""
        ;;
esac
