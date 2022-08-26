// Register plugin
if (hlive.afterMessage.get("hhhist") === undefined) {
    hlive.afterMessage.set("hhhist", function () {
        document.querySelectorAll("[__pushAttrName__]").forEach(function (el) {
            const path = el.getAttribute("__pushAttrName__")

            if (path !== "<USED>") {
                el.setAttribute("__pushAttrName__", "<USED>")

                history.pushState({path: path}, "", "__base_path__" + path)
            }
        });
    });

    // Init
    const p = (location.pathname+location.search).substring("__base_path__".length);
    history.replaceState({path: p}, null, p);

    onpopstate = function (event) {
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
}
