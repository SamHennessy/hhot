// Register plugin
hlive.afterMessage.push(function () {
    document.querySelectorAll("[__pushAttrName__]").forEach(function (el) {
        const path = el.getAttribute("__pushAttrName__")

        if (path !== "<USED>") {
            el.setAttribute("__pushAttrName__", "<USED>")

            history.pushState({path: path}, "", "__base_path__" + path)
        }
    });
});

window.onpopstate = function(event) {
    let path = "/"

    if (event.state && event.state.path) {
        path = event.state.path
    }

    hlive.sendMsg({
        t: "e",
        i: "__bindingID__",
        e: {
            "path": path
        }
    });
}
