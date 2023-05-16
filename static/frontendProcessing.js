function postProcessHTML() {
    assignBootstrapCSSClasses()
    registerKeyListener()
    appendCopyButtons()
}

function assignBootstrapCSSClasses() {
    const tables = document.querySelectorAll("table")
    tables.forEach(table => {
        table.setAttribute("class", "table")
    })

    const blockquotes = document.querySelectorAll("blockquote")
    blockquotes.forEach(blockquote => {
        blockquote.setAttribute("class", "text-secondary fst-italic")
    })
}

function registerKeyListener() {
    document.addEventListener("keyup", (e) => {
        if (e.key == "/") {
            document.getElementById("search-input-box").focus()
        }
    })
}

function appendCopyButtons() {
    const preBlocks = document.querySelectorAll("pre")
    preBlocks.forEach(block => {
        block.setAttribute("class", "d-flex")

        const codeBlocks = block.querySelectorAll("code")
        codeBlocks.forEach(codeBlock => {
            codeBlock.setAttribute("class", "col p-2 me-2 font-monospace bg-light border border-secondary rounded")
        })

        const copyButton = document.createElement("button")
        copyButton.innerHTML = "Copy"
        copyButton.setAttribute("type", "button")
        copyButton.setAttribute("class", "col-1 align-self-start btn btn-primary ")
        copyButton.addEventListener("click", handleCopyClick)
        block.append(copyButton)
    })
}

function handleCopyClick(evt) {
    const { children } = evt.target.parentElement
    const { innerText } = Array.from(children)[0]
    copyToClipboard(innerText)
}

const copyToClipboard = str => {
    const el = document.createElement("textarea")
    el.value = str
    el.setAttribute("readonly", "")
    el.style.position = "absolute"
    el.style.left = "-9999px"
    document.body.appendChild(el)
    const selected =
        document.getSelection().rangeCount > 0
            ? document.getSelection().getRangeAt(0)
            : false
    el.select()
    document.execCommand("copy")
    document.body.removeChild(el)
    if (selected) {
        document.getSelection().removeAllRanges()
        document.getSelection().addRange(selected)
    }
}
