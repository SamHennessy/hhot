function hhui_onloadurl(event, message) {
    if (event.type !== "load") {
        return message
    }



    console.log(event.target.contentWindow.location, message);

    if (event.target.contentWindow && event.target.contentWindow.location) {
        if (!message.e) {
            message.e = {};
        }

        message.e.location = event.target.contentWindow.location.href;
    }

    return message
}

// Register plugin
hlive.beforeSendEvent.push(hhui_onloadurl);
