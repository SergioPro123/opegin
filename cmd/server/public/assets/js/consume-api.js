$(document).on('submit', 'form#form-generate-doc', function (e) {
    e.preventDefault();
    $('.modal').modal('show');
    $('#bt-generate-doc').attr('disabled', 'disabled');
    let fechas = [
        'Enero',
        'Febrero',
        'Marzo',
        'Abril',
        'Mayo',
        'Junio',
        'Julio',
        'Agosto',
        'Septiembre',
        'Octubre',
        'Noviembre',
        'Diciembre',
    ];
    let dateDoc = $('#f1-date-doc').val().split('-');

    let sundayForm = {
        month: fechas[dateDoc[1] - 1],
        year: dateDoc[0],
        entry_time: $('#entry_time').val(),
        entry_time_sunday: $('#entry_time_sunday').val(),
        departure_time_sunday: $('#departure_time_sunday').val(),
        justification: 'Cumplimiento de meta del proyecto de despliegue  FTTH en el mes de DICIEMBRE.',
        responsible: {
            name: $('#responsible_name').val(),
            position: {
                name: $('#responsible_position').val(),
            },
        },
        immediate_boss: {
            name: $('#immediate_boss_name').val(),
            location: $('#immediate_boss_location').val(),
            department: $('#immediate_boss_department').val(),
        },
    };
    sundayFormString = JSON.stringify(sundayForm);

    var fileExcelEmployee = $('#file_employee')[0].files;

    var fd = new FormData();
    fd.append('sundayForm', sundayFormString);
    // Check file selected or not
    if (fileExcelEmployee.length > 0) {
        fd.append('file', fileExcelEmployee[0]);
    }
    $.ajax({
        url: '/api/v1/sundays',
        data: fd,
        cache: false,
        xhr: function () {
            var xhr = new XMLHttpRequest();
            xhr.onreadystatechange = function () {
                if (xhr.readyState == 2) {
                    if (xhr.status == 200) {
                        xhr.responseType = 'blob';
                    } else {
                        xhr.responseType = 'text';
                    }
                }
            };
            return xhr;
        },
        processData: false,
        contentType: false,
        type: 'POST',
        success: function (data, textStatus, request) {
            $('#bt-generate-doc').removeAttr('disabled');

            if (request.status === 200) {
                var blob = data;
                var filename = '';
                var disposition = request.getResponseHeader('Content-Disposition');
                if (disposition && disposition.indexOf('attachment') !== -1) {
                    var filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
                    var matches = filenameRegex.exec(disposition);
                    if (matches != null && matches[1]) filename = matches[1].replace(/['"]/g, '');
                }

                if (typeof window.navigator.msSaveBlob !== 'undefined') {
                    // IE workaround for "HTML7007: One or more blob URLs were revoked by closing the blob for which they were created. These URLs will no longer resolve as the data backing the URL has been freed."
                    window.navigator.msSaveBlob(blob, filename);
                } else {
                    var URL = window.URL || window.webkitURL;
                    var downloadUrl = URL.createObjectURL(blob);

                    if (filename) {
                        // use HTML5 a[download] attribute to specify filename
                        var a = document.createElement('a');
                        // safari doesn't support this yet
                        if (typeof a.download === 'undefined') {
                            window.location.href = downloadUrl;
                        } else {
                            a.href = downloadUrl;
                            a.download = filename;
                            document.body.appendChild(a);
                            a.click();
                        }
                    } else {
                        window.location.href = downloadUrl;
                    }

                    setTimeout(function () {
                        URL.revokeObjectURL(downloadUrl);
                    }, 100); // cleanup
                }

                $('.modal').modal('hide');
            }
        },
        error: function (error) {
            $('.modal').modal('hide');
            $('#bt-generate-doc').removeAttr('disabled');
            alert('Error: ' + (error?.responseJSON?.message || error));
        },
    });
});
