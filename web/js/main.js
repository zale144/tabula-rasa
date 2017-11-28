var cachedData = {};

$(document).ready(function() {
    loadTabs();
    var links = $('[rel="import"]');
    $.each(links, function(ind, link) {
        var content = link.import.querySelector('template').content;
        var clone = document.importNode(content, true);
        document.querySelector('#content').appendChild(clone);
    });
    $('.fragment').hide();
    $('#overview').show();
    $(window).on('hashchange', function (e) {
        render(decodeURI(window.location.hash));
    });

    if(!window.location.hash) {
        location.hash = '#home';
    } else {
        $(window).trigger('hashchange');
    }

    $('#confirm-delete').on('click','.btn-ok',function(e) {
        var $modalDiv = $(e.delegateTarget);
        var id = $(this).data('id');
        var obj = $('h1').text();
        var deleteUrl = "/rest/" + obj + "?id=" + id;

        $modalDiv.addClass('loading');
        $.ajax({
            url : deleteUrl,
            type : 'DELETE',
            success : function(result) {
                $modalDiv.modal('hide').removeClass('loading');
                cachedData[obj] = $.grep(cachedData[obj], function (o) {
                        return o.id == id || o[Object.keys(o)[0]] == id }, true);
                /*for (l in cachedData) {
                    var list = cachedData[l];
                    for(e in list) {
                        var keys = Object.keys(list[e]);
                        for(f in keys) {
                            if(obj.includes(keys[f]) && list[e][keys[f]].id == id) {
                                list[e][keys[f]] = undefined;
                            }
                        }
                    }
                }*/
                getData(obj, loadTable);
                loadTabs();
            }
        });
    });

    $('#confirm-delete').on('show.bs.modal',function(e) {
        var data = $(e.relatedTarget)
            .parent().parent()
            .children();
        var title = data.eq(2).text();
        var id = data.eq(1).text();
        if (data.length <= 4) {
            title = id;
        }
        $('.title', this).text("Are you sure you want to delete the item: \""
            + title
            + "\" from your collection?");
        $('.btn-ok', this).data('id', id);
    });

});

function formSubmit(e) {
    var idfld = $(e).find($('[datafld="id"]'));
    var id = idfld.val()?idfld.val():"";
    var inputs = $(e).find('input, select');
    var obj = $(e).attr("name");
    var data = {};

    for (var i = 0; i < inputs.length; i++) {
        var fld = $(inputs[i]).attr('datafld');
        if (fld) {
            var value = $(e).find($('[datafld="' + fld + '"]')).val();
            data[inputs[i].name] = value;
            if (obj == 'home') {
                data[inputs[i].name] = value.replace(/ /g, '_');
            }
        }
    }
    $.ajax({
        type: 'POST',
        url: "/rest/" + obj,
        data: JSON.stringify(data),
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        'dataType': 'json',
        success: function(newData, status) {
            if (newData) {
                id = newData.id;
            }
            cachedData[obj] = [];
            if (obj === 'home') {
                obj = data['Table name'];
                loadTabs();
                window.location.href = "#" + obj;
            }
            getData(obj, loadTable);
        }
    });
}

function render(url) {
    var obj = url.split('/')[0].substring(1);
    $("ul li").removeClass();
    $('[tab="' + obj + '"]').addClass('active');
    if (!obj.startsWith("add")) {
        getPage("overview");
        getData(obj, loadTable);
    } else {
        loadAddEditPage(obj);
    }
    $('h1').text(obj).css('textTransform', 'capitalize');
}

function getPage(obj) {
    $('.fragment').hide();
    $('#' + obj).show();
}

function getData(obj, func) {
    var data = cachedData[obj];
    if (data && data.length > 0) {
        func(data, obj)
    } else {
        loadAjax(obj, func);
    }
}

function loadTabs() {
    getData('home', dataF);
    function dataF(data) {
        var tabs = $('#tabs');
        tabs.empty();
        var name = 'home';
        var ul = $('<ul>').addClass('nav nav-tabs');
        var li = getLi(name);
        ul.append(li);
        tabs.append(ul);
        for (var i = 0; i < data.length; i++) {
            name = data[i].table_name;
            li = getLi(name);
            ul.append(li);
        }
        li = getLi('+');
        li.find('a').attr('href', '#addEdithome');
        ul.append(li);
        function getLi(name) {
            var li = $('<li>').attr('role', "presentation").attr('tab', name);
            var a = $('<a>').attr('href', '#' + ((name=='+')?'addEdit':name)).text(capitalizeFirstLetter(name));
            return li.append(a);
        }
    }
}

function loadTable(data, obj) {
    var addButton = $('#add').text('Add '
        + capitalizeFirstLetter(((obj=='home')?'Entity':obj)));
    if (obj === 'home') {
        addButton.removeAttr('onclick role');
        addButton.attr('href', '#addEdithome');
    } else {
        addButton.removeAttr('href');
        addButton.attr('onclick', 'addRow("' + obj + '")').attr('role', 'button');
    }
    var table = $('<table></table>').addClass('table table-bordered');
    var grid = $('#grid');
    grid.empty();
    var headers = $('<thead></thead>');
    var contentRows = $('<tbody></tbody>');

    grid.append(table);
    table.append(contentRows);
    table.append(headers);

    if (data.length > 0) {
        var keys = Object.keys(data[0]);
        loadTableHeaders(headers, keys);
        $.each(data, function (index, item) {
            var row = $('<tr></tr>');
            contentRows.append(row);
            row.append($('<input type="hidden" value="' + item.id + '" />'));
            for (var k in keys) {
                var td = $('<td></td>');
                row.append(td);
                if (item[keys[k]]) {
                    var str;
                    if (item[keys[k]] instanceof  Object &&
                        !(item[keys[k]] instanceof  Array)) {
                        str = item[keys[k]].name;
                    } else if (keys[k] == 'date') {
                        str = $('<a href="' + "" + item.id + '" >' + formatDate(item[keys[k]]) + '</a>');
                    } else if (item[keys[k]] instanceof  Array) {
                        str = $('<ul></ul>');
                        item[keys[k]].forEach(function(e) {
                            str.append($('<li>' + e.name + '</li>'));
                        });
                    } else {
                        str = (k == 1)?$('<a href="' + "" + item.id + '" >' + item[keys[k]] + '</a>'): item[keys[k]];
                    }
                    td.append(str);
                }

            }
            row.append($('<td><a role="button" class="btn btn-info" onclick="editRow(this.parentElement.parentElement)">Edit</a></td>'));
            row.append($('<td><a role="button" class="btn btn-danger" data-toggle="modal" ' +
                ' data-target="#confirm-delete">Delete</a></td>'));
        });
    }

    function formatDate(date) {
        var d = new Date(date);
        d = ("0" + d.getDate()).slice(-2) + "-" + ("0"+(d.getMonth()+1)).slice(-2) + "-" + d.getFullYear() + " "
            + ("0" + d.getHours()).slice(-2) + ":" + ("0" + d.getMinutes()).slice(-2) + ":" + ("0" + d.getSeconds()).slice(-2);
        return d;
    }
}

function loadTableHeaders(headers, keys) {
    var tH = $('<tr></tr>');
    headers.append(tH);
    for(k in keys) {
        tH.append($('<th>' + capitalizeFirstLetter(keys[k]) + '</th>'));
    }
    tH.append($('<th>Edit</th>')).append($('<th>Delete</th>'));
}

function loadAjax(obj, func) {
    $.ajax({
        url : "/rest/" + obj,
        type : 'GET',
        success : function(data, status) {
            var jsonData = JSON.parse(data);
            cachedData[obj] = jsonData;
            func(jsonData, obj);
        },
        error: function() {
            $('#errorMessages')
                .append($('<li>')
                    .attr({class: 'list-group-item list-group-item-danger'})
                    .text('Error calling web service.  Please try again later.'));
        }
    });
}

function loadAddEditPage(obj, id) {
    dropdown();
    obj = obj.substring(7, obj.length);
    getPage('addEdit');

    var formDiv = $('#form');
    formDiv.empty();
    var form = $('<form>').addClass('form-horizontal').attr('name', obj)
        .attr('onsubmit', 'event.preventDefault();formSubmit(this);');
    formDiv.append(form);

    if (obj !== 'home') {
        getData(obj + '/cols', function(data) {
            var cols = [];
            for (var d in data) {
                cols.push(data[d].COLUMN_NAME);
            }
            createFields(obj, cols);
        });
    } else {
        createFields(obj, ['Table name', 'Column name']);
    }

    function createFields(obj, cols) {
        for (var i = 0; i < cols.length; i++) {
            var col = cols[i];
            if (col === 'id') {
                continue;
            }
            var div = $('<div>').addClass('form-group');
            form.append(div);
            var label = $('<label>').addClass('col-md-4 control-label')
                .attr('for', col).text(capitalizeFirstLetter(col) + ':');
            div.append(label);
            var fDiv = $('<div>').addClass('col-md-8');
            div.append(fDiv);
            var input = $('<input>').attr('type', 'text').attr('name', col)
                .attr('onclick', 'this.select();');
            fDiv.append(input);
            if (col === 'Column name') {
                var ol = $('<ol>').addClass('list-group col-md-offset-2').attr('id', 'addedCols');
                div.after(ol);
                var add = $('<button>').attr('onclick', 'event.preventDefault();addColumn(this)').text('+');
                fDiv.append(add);
            } else {
                input.attr('datafld', col);
            }
        }
        var buttonDiv = $('<div>').addClass('form-group');
        form.append(buttonDiv);
        var cancelDiv = $('<div>').addClass('col-md-offset-2 col-md-2');
        buttonDiv.append(cancelDiv);
        var a = $('<a>').addClass('btn btn-default').attr('href', '#' + obj).text('Cancel');
        cancelDiv.append(a);
        var submitDiv = $('<div>').addClass('col-md-2');
        buttonDiv.append(submitDiv);
        var submit = $('<input>').addClass('btn btn-default').attr('type', 'submit')
            .attr('value', 'Create ' + capitalizeFirstLetter(obj));
        submitDiv.append(submit);
    }
}

function addColumn(e) {
    var input = $(e).prev();
    var name = input.val();
    if (name) {
        var addedCols = $('#addedCols');
        var li = $('<li>').addClass('list-group-item');
        var fIn = $('<input>').val(name).attr('datafld', name).attr('name', name).prop('required', true);
        li.append(fIn);
        addedCols.append(li);
        input.val('');
    } else {
        alert('Please input column name.');
    }
}

function addRow(obj) {
    getData(obj + '/cols', function(data) {
        var headers = $('thead');
        var cols = [];
        for (var d in data) {
            cols.push(data[d].COLUMN_NAME);
        }
        if ($('thead').children().length === 0) {
            loadTableHeaders(headers, cols);
        }
        var body = $('#grid').find('tbody');
        var form = $('<form>').attr('id', 'form').attr('name', obj)
            .attr('onsubmit', 'event.preventDefault();formSubmit(this);');
        var tr = $('<tr>').attr('name', obj);
        form.append(tr);
        body.append(tr);
        for (var i = 0; i < cols.length; i++) {
            var td = $('<td>');
            tr.append(td);
            if (i > 0) {
                var input = $('<input>').attr('type', 'text').attr('name', cols[i])
                    .attr('onclick', 'this.select();').attr('datafld', cols[i])
                    .attr('placeholder', capitalizeFirstLetter(cols[i]));
                td.append(input);
            }
        }
        tr.append($('<td><input type="submit" class="btn btn-success" ' +
            'onclick="formSubmit(this.parentElement.parentElement)" /></td>'));
        tr.append($('<td><a role="button" class="btn btn-default" ' +
            'onclick="removeRow(this.parentElement.parentElement);">Cancel</a></td>'));

    });
}

function removeRow(e) {
    $(e).remove();
}

function editRow(e) {
    var row = $(e);
    var obj = window.location.hash.substring(1);
    row.attr('id', 'form').attr('name', obj);
    var tds = row.children('td');
    var relFields = tds.slice(0, tds.length - 2);
    var headers = $('thead').find('th');
    relFields.each(function(i, el) {
        var value = $(el).text();
        $(el).text('');
        var col = capitalizeFirstLetter($(headers[i]).text());
        var input = $('<input>').val(value).attr('type', 'text').attr('name', col)
            .attr('onclick', 'this.select();').attr('datafld', col)
        if (i == 0) {
            input.prop('disabled', true);
            input.width('100%');
            $(el).width('30px');
        }
        $(el).append(input);
    });
    tds.last().find('a').removeClass('btn-danger').addClass('btn-default')
        .text('Cancel').removeAttr('data-target data-toggle').click(function() {getData(obj, loadTable)});
    tds.last().prev().find('a').removeClass('btn-info').addClass('btn-success').text('Save')
        .attr('onclick', 'formSubmit(this.parentElement.parentElement)');
}

function dropdown() {
    elems = $('select');
    elems.empty();
    $.each(elems, function(i, elem) {
        var obj = $(elem).attr('datafld');
        getData(obj, function(data) {
            $.each(data, function(ind, d) {
                var option = document.createElement('option');
                option.innerHTML = d.name;
                option.value = d.id;
                elem.append(option);
            });
        });
    });
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}

