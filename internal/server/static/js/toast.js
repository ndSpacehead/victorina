class Toast {
    /**
     * A class representing a Toast notification.
     * @param level {("info"|"success"|"warning"|"danger")}
     * @param message { string } base64 string
     */
    constructor(level, message) {
        this.level = level;
        this.message = base64ToString(message);
    }

    /**
     * Makes the toast container element. A button containing the entire notification.
     * @returns {HTMLButtonElement}
     */
    #makeToastContainerButton() {
        const button = document.createElement("button");
        button.classList.add("alert");
        button.classList.add(`alert-${this.level}`);
        button.classList.add("d-flex");
        button.classList.add("align-items-center");
        button.setAttribute("role", "alert");
        button.setAttribute("aria-label", "Close");
        button.addEventListener("click", () => button.remove());
        return button;
    }

    /**
     * Makes the element containing the body of the toast notification.
     * @returns {HTMLSpanElement}
     */
    #makeToastContentElement() {
        const messageContainer = document.createElement("span");
        messageContainer.textContent = this.message;
        return messageContainer;
    }

    /**
     * Presents the toast notification at the end of the given container.
     * @param containerQuerySelector {string} a CSS query selector identifying the container for all toasts.
     */
    show(containerQuerySelector = "#toast-container") {
        const toast = this.#makeToastContainerButton();
        const toastContent = this.#makeToastContentElement()
        toast.appendChild(toastContent);

        const toastContainer = document.querySelector(containerQuerySelector);
        toastContainer.appendChild(toast);

        setTimeout(() => {
            toast.remove();
        }, 5000);
    }
}

/**
 * Presents a toast notification when the `makeToast` event is triggered.
 * @param e {{detail: {level: string, message: string}}}
 */
function onMakeToast(e) {
    console.log(e);
    const toast = new Toast(e.detail.level, e.detail.message);
    toast.show();
}

/**
 * Decodes base64 to UTF-8 string.
 * @param encoded {string} base64 string
 * @returns {string} decoded UTF-8 string
 */
function base64ToString(encoded) {
    return decodeURIComponent(atob(encoded).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
}

document.body.addEventListener("makeToast", onMakeToast);