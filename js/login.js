$(document).ready(function () {
    var error = '{{.error}}';
    if (error) {
        $('#errorMessage').text(error);
        $('#errorModal').modal('show');
    }
});