console.log('Admin panel loaded');

function adminSyncColorToText(colorInput) {
    var textId = colorInput.id.replace('-picker', '');
    var textInput = document.getElementById(textId);
    if (textInput) {
        textInput.value = colorInput.value.toUpperCase();
    }
}

function adminSyncTextToPicker(textInput) {
    var picker = document.getElementById(textInput.id + '-picker');
    if (picker && /^#[0-9a-fA-F]{6}$/.test(textInput.value)) {
        picker.value = textInput.value;
    }
}
