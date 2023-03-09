class Justification {
    constructor(id, month, description) {
        this.Id = id;
        this.Month = month;
        this.Description = description;
    }
}
var justifications = [];

$(document).ready(() => {
    justifications.push(new Justification(1, `justification_month_1`, 'justification_description_1'));
});

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
        justification: [],
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
    //add justifications
    justifications.every((justification) => {
        let justification_month = $('#' + justification.Month);
        let justification_description = $('#' + justification.Description);
        if (!(justification_month.length && justification_description.length)) {
            return true;
        }
        description = justification_description.val();
        month = parseInt(justification_month.val());
        if (isNaN(month)) {
            $('.modal').modal('hide');
            $('#bt-generate-doc').removeAttr('disabled');
            alert('Seleccione un mes en la justificación.');
            return false;
        }
        sundayForm.justification.push({
            description: description,
            number_month: month,
        });
        return true;
    });
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

function deleteFormGroupJustify(id) {
    let i = 0;
    justifications.every((justification) => {
        if (justification.Id != id) {
            i++;
            return true;
        }
        justifications.splice(i, 1);
        return false;
    });
    $('#form-grup-justify-' + id).remove();
}

$(document).on('click', '#add_justification', function () {
    if (justifications.length >= 12) {
        alert('No se pueden agregar mas de 12 justificaciones, debido a que solo se puede 1 por mes.');
        return;
    }
    let newId = justifications[justifications.length - 1].Id + 1;
    let justification = new Justification(newId, `justification_month_${newId}`, `justification_description_${newId}`);
    justifications.push(justification);

    let html =
        `<div class='form-group' id="form-grup-justify-` +
        justification.Id +
        `">
            <div class='form-group col-md-3'>
                <label class='' for='` +
        justification.Month +
        `'>Mes</label>
                <select
                    name='` +
        justification.Month +
        `'
                    id='` +
        justification.Month +
        `'
                    class='form-control'
                >
                    <option value=''>--Selecione el mes--</option>
                    <option value='01'>Enero</option>
                    <option value='02'>Febrero</option>
                    <option value='03'>Marzo</option>
                    <option value='04'>Abril</option>
                    <option value='05'>Mayo</option>
                    <option value='06'>Junio</option>
                    <option value='07'>Julio</option>
                    <option value='08'>Agosto</option>
                    <option value='09'>Septiembre</option>
                    <option value='10'>Octubre</option>
                    <option value='11'>Noviembre</option>
                    <option value='12'>Diciembre</option>
                </select>
            </div>
            <div class='form-group col-md-8'>
                <label class='' for='` +
        justification.Description +
        `'>Descripción</label>
                <textarea
                required
                    class='form-control'
                    id='` +
        justification.Description +
        `'
                    name='` +
        justification.Description +
        `'
                    rows='3'
                ></textarea>
            </div>

            <div class='form-group col-md-1 justify-content-md-center'>
                <button onClick="deleteFormGroupJustify(` +
        justification.Id +
        `)" type='button' class='close form-control' aria-label='Close'>
                    <span aria-hidden='true'>&times;</span>
                </button>
    </div>
</div>
    `;
    $('#new_justification').append(html);
});
