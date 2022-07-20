// Iframe Watcher
// If already running don't start a new loop
// Will self terminate if attribute is removed
if (!window.hhotIframeWatcherDelay) {
    window.hhotIframeWatcherDelay = __iframeWatcherDelay__

    watcher = function () {
        const els = document.querySelectorAll("[data-hlive-on*=__iframeWatcherEvent__]")

        if (els.length === 0) {
            window.hhotIframeWatcherDelay = false

            return
        }

        els.forEach(function (el) {
            const ids = hlive.getEventHandlerIDs(el)

            // Change?
            // TODO: el.contentDocument can be null
            if (el.hhotIframeWatcherTitle === el.contentDocument.title &&
                el.hhotIframeWatcherPathname === el.contentWindow.location.pathname
            ) {
                return
            }

            el.hhotIframeWatcherTitle = el.contentDocument.title
            el.hhotIframeWatcherPathname = el.contentWindow.location.pathname

            if (!ids["__iframeWatcherEvent__"] || ids["__iframeWatcherEvent__"].length < 1) {
                return
            }

            for (let i = 0; i < ids["__iframeWatcherEvent__"].length; i++) {
                hlive.sendMsg({
                    t: "e",
                    i: ids["__iframeWatcherEvent__"][i],
                    e: {
                        "title": el.hhotIframeWatcherTitle,
                        "path": el.hhotIframeWatcherPathname,
                    }
                });
            }
        });

        if (window.hhotIframeWatcherDelay) {
            setTimeout(watcher, window.hhotIframeWatcherDelay)
        }
    }

    // TODO: we are getting very long setTimeout chains
    setTimeout(watcher, window.hhotIframeWatcherDelay)
} else {
    // Delay may have changed between page loads
    window.hhotIframeWatcherDelay = __iframeWatcherDelay__
}
