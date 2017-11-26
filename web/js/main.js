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

    for (var i = 0; i < inputs.length - 1; i++) {
        var fld = $(inputs[i]).attr('datafld');
        if (fld) {
            data[inputs[i].name] = $(e).find($('[datafld="' + fld + '"]')).val();
        }
    }
    $.ajax({
        type: id?'PUT':'POST',
        url: "/rest/" + obj + "/" + id,
        data: JSON.stringify(data),
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        'dataType': 'json',
        success: function(data, status) {
            if (data) {
                id = data.id;
            }
            cachedData[obj] = [];
            window.location.href = "#" + obj;
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
        ul.append(li);
        function getLi(name) {
            var li = $('<li>').attr('role', "presentation").attr('tab', name);
            var a = $('<a>').attr('href', '#' + ((name=='+')?'addEdit':name)).text(capitalizeFirstLetter(name));
            return li.append(a);
        }
    }
}

function loadTable(data, obj) {
    $('#add').attr('href', '#addEdit' + obj).text('Add '
        + capitalizeFirstLetter(((obj=='home')?'Entity':obj)));
    var grid = $('#grid');
    grid.empty();
    if (data.length > 0) {
        var table = $('<table></table>').addClass('table table-bordered');
        var headers = $('<thead></thead>');
        var contentRows = $('<tbody></tbody>');
        var keys = Object.keys(data[0]);

        table.append(headers);
        table.append(contentRows);
        grid.append(table);

        var tH = $('<tr></tr>');
        headers.append(tH);
        for(k in keys) {
            tH.append($('<th>' + capitalizeFirstLetter(keys[k]) + '</th>'));
        }
        tH.append($('<th>Edit</th>')).append($('<th>Delete</th>'));
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
            //row += '<td><a class="btn btn-info" role="button" onclick="loadAddEditPage(\''+obj+'\', '+item.id+')">Edit</a></td>';
            row.append($('<td><a class="btn btn-info" href="#' + '' + '">Edit</a></td>'));
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

    var cols = [];
    getData(obj + '/cols', function(data) {
       for (var d in data) {
           cols.push(data[d].COLUMN_NAME);
       }
        createFields(cols);
    });

    function createFields(cols) {
        var formDiv = $('#form');
        formDiv.empty();
        var form = $('<form>').addClass('form-horizontal').attr('name', obj)
            .attr('onsubmit', 'event.preventDefault();formSubmit(this);');
        formDiv.append(form);
        for (var i = 0; i < cols.length; i++) {
            var col = cols[i];
            if (col == 'id') {
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
                .attr('datafld', col);
            fDiv.append(input);
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

