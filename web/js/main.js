var cachedData = {'types':[
        {id : 'varchar(45)', name : 'Text'},
        {id : 'int(11)',         name : 'Integer'},
        {id : 'double',       name : 'Decimal'},
        {id : 'REF',         name : 'Reference'}]};

$(document).ready(function() {
    (function($) {
        $.fn.changeElementType = function(newType) {
            var attrs = {};

            $.each(this[0].attributes, function(idx, attr) {
                attrs[attr.nodeName] = attr.nodeValue;
            });

            this.replaceWith(function() {
                return $("<" + newType + "/>", attrs).append($(this).contents());
            });
        };
    })(jQuery);
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

    $(document).on('keydown', function (e) {
        if (e.which === 13) {
            e.preventDefault();
            //this.blur();
            window.focus();
            $('.btn-success').click();
            formProcess(this.activeElement.parentElement);
            return false;
        }
    });

    $(function() {
        /*$('*:not(input,select,td)').click(function() {
            if ($('.emptyDropdown').length > 0) {
                formProcess();
            }
            return false;
        });*/

    });

});

function saveAjax(data, obj, after) {
    $.ajax({
        type: 'POST',
        url: "/rest/" + obj,
        data: JSON.stringify(data),
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        'dataType': 'json',
        success: function(newData) {
            after();
        }
    });
}

function formProcess(e) {
    //var caller = $(e.parentElement);
    var inputs = $('input, select');
    var obj = $('table, form').attr('name');
    var data = {};

    for (var i = 0; i < inputs.length; i++) {
        var fld = $(inputs[i]).attr('datafld');
        if (fld) {
            var attr = /*caller.find(*/$('[datafld="' + fld + '"]')/*)*/.val();
            if (attr && obj === 'home') {
                var value = /*caller.find(*/$('select[datafld="' + fld + '"][name="types"]')/*)*/.val();
                attr = attr.replace(/ /g, '_');
                if (i === 0) {
                    data[inputs[i].name] = attr;
                } else {
                    if (value === 'REF') {
                        var select =/*caller.find(*/$('select[datafld="' + fld + '"][name="home"]')/*)*/.val();
                        data[fld] = select + '_id';
                    } else {
                        data[attr] = value;
                    }
                }
            } else {
                var key = $(inputs[i]).attr('datafld');
                data[key] = attr;
            }
        }
    }
    saveAjax(data, obj, function () {
        cachedData[obj] = [];
        if (obj === 'home') {
            obj = data['Table name'];
            loadTabs();
            window.location.href = "#" + obj;
        }
        getData(obj, loadTable);
    });
}

function render(url) {
    var obj = url.split('/')[0].substring(1);
    $("ul li").removeClass();
    $('[tab="' + obj + '"]').addClass('active');
    if (!obj.startsWith("add")) {
        getPage("overview");
        getData(obj, loadTable);
        $('h1').text(obj).css('textTransform', 'capitalize');
    } else {
        if (!$('[name="Table name"]')[0]) {
            loadAddEditPage(obj);
        } else {
            getPage('addEdit');
        }
        $('h1').text("Add Table").css('textTransform', 'capitalize');
    }
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
            var a = $('<a>').attr('href', '#' + ((name==='+')?'addEdit':name)).text(capitalizeFirstLetter(name));
            return li.append(a);
        }
    }
}

function loadTable(data, obj) {
    var addButton = $('#add').text('Add '
        + capitalizeFirstLetter(((obj==='home')?'Table':obj)));
    if (obj === 'home') {
        addButton.removeAttr('onclick role');
        addButton.attr('href', '#addEdithome');
    } else {
        addButton.removeAttr('href');
        addButton.attr('onclick', 'addRow("' + obj + '")').attr('role', 'button');
    }
    var table = $('<table>').addClass('table table-bordered').attr('name', obj);
    var grid = $('#grid');
    grid.empty();
    var headers = $('<thead>');
    var contentRows = $('<tbody>');

    grid.append(table);
    table.append(contentRows);
    table.append(headers);

    if (data.length > 0) {
        var keys = Object.keys(data[0]);
        loadTableHeaders(headers, keys);
        $.each(data, function (index, item) {
            var row = $('<tr>');
            contentRows.append(row);
            row.append($('<input type="hidden" value="' + item.id + '" />').attr('datafld', 'id'));
            for (var k = 0; k < keys.length; k++) {
                var td = $('<td>');
                if (k === 0) {
                    td.addClass('col-md-1');
                } else {
                    td.addClass('col-md-2').attr('onclick', 'editCell(this);');
                }
                row.append(td);
                if (item[keys[k]]) {
                    var str;
                    if (item[keys[k]] instanceof  Object &&
                        !(item[keys[k]] instanceof  Array)) {
                        str = $('<div>').text(item[keys[k]].name)
                    } else if (keys[k] === 'date') {
                        str = $('<a href="' + "" + item.id + '" >' + formatDate(item[keys[k]]) + '</a>');
                    } else if (item[keys[k]] instanceof  Array) {
                        str = $('<ul>');
                        item[keys[k]].forEach(function(e) {
                            str.append($('<li>' + e.name + '</li>'));
                        });
                    } else {
                        str = (k === 1)?$('<a href="' + "" + item.id + '" >' + item[keys[k]] + '</a>'): $('<div>')
                            .text(item[keys[k]]);
                    }
                    td.append(str);
                }

            }
            var etd = $('<td>').addClass('col-md-1');
            var editBtn = $('<a role="button" class="btn" >');
            etd.append(editBtn);
            row.append(etd);
            if (obj === 'home') {
                editBtn.addClass('btn-default').attr('onclick', 'window.location.href = "#" + $(":nth-child(2)", ' +
                    'this.parentElement.parentElement).text()')
                    .text('View');
            } else {
                editBtn.addClass('btn-info').attr('onclick', 'editRow(this.parentElement.parentElement)').text('Edit');
            }
            row.append($('<td><a role="button" class="btn btn-danger" data-toggle="modal" ' +
                ' data-target="#confirm-delete">Delete</a></td>').addClass('col-md-1'));
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
    for(var k = 0; k < keys.length; k++) {
        var key = keys[k];
        if (key.endsWith("_fk")) {
            key = key.substring(0, key.length-3);
        }
        key = key.replace(/_/g, ' ');
        tH.append($('<th>' + capitalizeFirstLetter(key) + '</th>'));
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
        }
    });
}

function loadAddEditPage(obj, id) {
    obj = obj.substring(7, obj.length);
    getPage('addEdit');

    var formDiv = $('#form');
    formDiv.empty();
    var form = $('<form>').addClass('form-horizontal').attr('name', obj)
        .attr('onsubmit', 'event.preventDefault();checkColumns(this);');
    formDiv.append(form);

    if (obj !== 'home') {
        getData(obj + '/cols', function(data) {
            var cols = [];
            for (var d in data) {
                cols.push(data[d].Field);
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
            input = $('<input>').attr('type', 'text').attr('name', col)
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
            .attr('value', 'Create table');
        submitDiv.append(submit);
    }
    $('input[name="Table name"]').on('input',function(e){
        if (e.currentTarget.value === "") {
            $('input[type="submit"].btn-default').get(0).value = "Create table"
        } else {
            $('input[type="submit"].btn-default').get(0).value = "Create table '" + e.currentTarget.value + "'";
        }
    });
}

function addColumn(e) {
    var input = $(e).prev();
    var name = input.val();
    if (name) {
        var addedCols = $('#addedCols');
        var li = $('<li>').addClass('list-group-item');
        var fIn = $('<input>').val(name).attr('datafld', name).attr('name', name).prop('required', true);
        var label = $('<label>Type</label>').attr('for', name);
        li.append(fIn);
        li.append(label);
        dropdown('types', name, li);
        addedCols.append(li);
        input.val('');
        $('select[name="types"][datafld="' + name + '"]').change(function(e) {
            if(this.value === 'REF') {
                fIn.siblings('label').remove();
                dropdown('home', name, li);
            } else {
                $(this).before(label);
                li.find('select[name="home"]').remove();
            }
        });
    } else {
        alert('Please input column name.');
    }
}

function addRow(obj) {
    getData(obj + '/cols', function(data) {
        var body = $('#grid').find('tbody');
        var form = $('<form>').attr('id', 'form').attr('name', obj)

        tr = $('<tr>').attr('name', obj);
        form.append(tr);
        body.append(tr);
        var cols = [];
        openRowInputs(tr, obj, data, cols);

        var headers = $('thead');
        if (headers.children().length === 0) {
            loadTableHeaders(headers, cols);
        }
        tr.append($('<td><input type="submit" class="btn btn-success" ' +
            'onclick="formProcess(this.parentElement)" /></td>'));
        tr.append($('<td><a role="button" class="btn btn-default" ' +
            'onclick="$(this.parentElement.parentElement).remove();">Cancel</a></td>'));
        if ($("input:text:visible:not(:disabled):first").get(0)) {
            $("input:text:visible:not(:disabled):first").get(0).focus();
        }
    });
}

function openRowInputs(tr, obj, data, cols, rowData) {
    var types = [];
    getColsAndTypes(data, cols, types);
    for (var i = 0; i < cols.length; i++) {
        var inputVal;
        var td = $('<td>').addClass('col-md-1');
        tr.append(td);
        if (rowData) {
            inputVal = rowData[i];
            td.addClass('col-md-2');
        }
        if (rowData || i > 0) {
            openCellInput(td, cols[i], types[i], inputVal);
        }
    }
}

function openCellInput(td, col, type, inputVal) {
    td.empty();
    if (type.REF) {
        dropdown(type.REF, col, td, inputVal);
    } else {
        var input = $('<input>').attr('name', col)
            .attr('onclick', 'this.select();').attr('datafld', col)
            .attr('placeholder', capitalizeFirstLetter(col)).val(inputVal);
        if ($('[datafld]:not([type=hidden])').length == 0) {
            input.attr('onblur', 'reloadTable();');
        }

        td.append(input);

        if ($.inArray(type, ['Integer', 'Decimal']) !== -1) {
            input.attr('type', 'number').val(input.val()?input.val():0);
            if (type === 'Decimal') {
                input.attr('step', '0.01');
            }
        }
    }
}

function editCell(e) {
    var td = $(e);
    td.removeAttr('onclick');
    var value = td.text();
    var obj = td.parent().parent().parent().attr('name');
    td.parent().attr('name', obj);
    var types = [];
    var cols = [];
    getData(obj + '/cols', function(data) {
        getColsAndTypes(data, cols, types);
        var col = cols[td.index()-1];
        var type = types[td.index()-1];
        openCellInput(td, col, type, value);
        $("input:visible:not(:disabled,:submit)").focus();
    });
}

function reloadTable() {
    var obj = window.location.hash.substring(1);
    getData(obj, loadTable);
}

function getColsAndTypes(data, cols, types) {
    for (var d = 0; d < data.length; d++) {
        cols.push(data[d].column_name);
        if (data[d].referenced_table_name !== '') {
            types.push({'REF': data[d].referenced_table_name});
        } else {
            var type = $.grep(cachedData.types, function(e){
                return e.id === data[d].column_type; })[0];
            types.push(type.name);
        }
    }
}

function checkColumns(e) {
    if (!$('[name="Table name"]').val()) {
        alert('Please input table name before saving');
    } else if ($('[name="Column name"]').val()) {
        alert('Please add the column before saving');
    } else if ($('#addedCols').is(':empty')) {
        alert('Please add at lease one column before saving');
    } else {
        formProcess(e);
    }
}

function editRow(e) {
    var obj = $(e).parent().parent().attr('name');
    getData(obj + '/cols', function(data) {

        var row = $(e);
        var rowData = row.children("td").map(function() {
            return $(this).text();
        }).get();

        var obj = window.location.hash.substring(1);
        row.attr('id', 'form').attr('name', obj);
        row.children().slice(2, row.children().length-2).remove();
        var cols = [];
        tr = $('<tr>').attr('name', obj);
        row.parent().append(tr);
        openRowInputs(tr, obj, data, cols, rowData);
        $.each(tr.children(), function(i, td) {
            if (i > 0) {
                row.children().eq(i).after(td);
            }
        });
        tr.remove();
        row.children().last().find('a').removeClass('btn-danger').addClass('btn-default')
            .text('Cancel').removeAttr('data-target data-toggle').click(function() {getData(obj, loadTable)});
        row.children().last().prev().find('a').removeClass('btn-info').addClass('btn-success').text('Save')
            .attr('onclick', 'formProcess(this.parentElement)');
        $("input:text:visible:not(:disabled):first").focus();
    });
}

function dropdown(name, datafld, parent, value) {
    var select = $('<select>').addClass('emptyDropdown').attr('name', name).attr('datafld', datafld);
    parent.append(select)
    getData(name, function(data) {
        $.each(data, function(ind, d) {
            var option = document.createElement('option');
            option.innerHTML = (name==='home')?d.table_name:d.name;
            option.value = (name==='home')?d.table_name:d.id;
            select.append(option);
        });
        $('.emptyDropdown').attr('onblur', 'reloadTable();').focus();
        $('select option:contains("' + value + '")').attr("selected","selected");
    });
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}

