function _goryuk_set_return() { return "${1:-0}"; }
goryuk_precmd() {
    echo "Go-Ryuk is running!!!"
    if [[ -z "${PROMPT_COMMAND-}" ]]; then
        PROMPT_COMMAND="go-ryuk check"
    elif [[ "$PROMPT_COMMAND" != *"go-ryuk check"* ]]; then
        GORYUK_PROMPT_COMMAND="$PROMPT_COMMAND"
        PROMPT_COMMAND="go-ryuk check"
    fi
}
