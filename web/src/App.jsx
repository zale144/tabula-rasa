const {Button, Modal} = ReactBootstrap;
// main component
class App extends React.Component {
    constructor(props) {
        super(props);
        this.selectRow = this.selectRow.bind(this);
        this.closeCellHandler = this.closeCellHandler.bind(this);
        this.state = {
            tableName: 'tables',
            modalOpen: false,
            isEditMode: false
        }
    }
    // method to change the name of the table when navigating between tabs
    changeTable(nextTable) {
        this.setState({
            tableName: nextTable,
        });
    }
    // method for loading the tabs
    loadTabs() {
        fetch("/rest/tables/tabs")
            .then(res => res.json(), {
                method: 'GET'
            })
            .then(
                (result) => {
                    const tabs = [{Table_name: 'tables'}, ...result];
                    this.setState({
                        tabsLoaded: true,
                        tabs: tabs,
                        tabsError: null
                    });
                },
                (error) => {
                    this.setState({
                        tabsLoaded: true,
                        tabsError: error,
                    });
                }
            )
    }
    // method for loading the table data
    loadTable() {
        fetch("/rest/" + this.state.tableName + "/rows", {
            method: 'GET'
        })
            .then(res => res.json()
            .then(r => r.map((v) => {
                return v.Table_name ? {"Table": v.Table_name} : v
            })))
            .then(
                (result) => {
                    if (this.state.intervalId) {
                        window.clearInterval(this.state.intervalId);
                    }
                    this.setState({
                        tableLoaded: true,
                        tableData: result,
                        tableError: null,
                        intervalId: null
                    });
                },
                (error) => {
                    this.setState({
                        tableLoaded: true,
                        tableError: error,
                    });
                    if (!this.state.intervalId) {
                        this.setState({
                            intervalId: setInterval(() => {
                                this.componentDidMount()
                            }, 1000)
                        });

                    }
                }
            )
    }
    // method for loading column names
    loadColumns() {
        fetch("/rest/" + this.state.tableName + "/cols", {
            method: 'GET'
        })
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        columns: result,
                        originalColLength: result.length,
                        columnsLoaded: true,
                        columnsError: null
                    });
                },
                (error) => {
                    this.setState({
                        columnsLoaded: true,
                        columnsError: error
                    });
                }
            )
    }
    // method for handling saving/updating cells, rows and tables
    saveTable(data, typ, afterSave) {
        fetch("/rest/" + this.state.tableName + "/" + typ, {
            method: 'POST',
            body: JSON.stringify(data)
        })
            .then(res => {
                return res.json()
            })
            .then(result => {
                afterSave(result);
            },
            (error) => {
                this.componentDidMount();
            });
    }
    // method for handling when user focuses out of a cell
    closeCellHandler(event) {
        const tableData = this.state.isEditMode ?
            this.state.columns.slice() : this.state.tableData.slice(),
            {value, name, attributes} = event.target,
            index = Number(attributes.rowindex.value);

        const oldName = index || this.state.originalColLength === tableData.length ?
            tableData[index]["Column_name"] : "";

        if (tableData[index][name] === value)
            return;

        if ([oldName, value].includes("id")) {
            alert("Changing the 'id' is not allowed!");
            return;
        }

        tableData[index][name] = value;
        let typ = "rows";
        let cellObj = {};

        if (this.state.isEditMode) {
            cellObj = tableData[index];
            cellObj["Old_column_name"] = oldName;
            typ = "cols";
        } else {
            cellObj.Id = tableData[index].Id;
            cellObj[name] = value;
        }
        const tempData = JSON.parse(JSON.stringify(tableData));
        tempData.map((v, i) => {
            if (i === index && v[name] === value) {
                v[name] = <img src="/static/public/ajax-loader.gif" style={{margin:'auto'}}/>;
            }
            return v;
        });
        this.setState({
            tableData: tempData
        });
        this.saveTable(cellObj, typ, (result) => {
            if (result[0].Id) {
                if (!tableData[index].Id) {
                    tableData[index].Id = result[0].Id;
                }
                this.setState({
                    tableData: tableData,
                    tableLoaded: true
                });
            } else if (tableData[index].Table) {
                tableData[index].Table = result[0].Table;
            }
            this.componentDidMount();
        });
    }
    // handling when user deletes a row or table
    handleDelete() {
        const tableName = this.state.tableName;
        const id = this.state.selectedRowId;
        const typ = this.state.isEditMode ? "cols" : "rows";

        fetch("/rest/" + tableName + "/" + typ + "?id=" + id, {
            method: 'DELETE'
        })
            .then(
                (result) => {
                    this.componentDidMount();
                    this.setState({
                        modalOpen: false,
                        selectedRowId: null
                    });
                },
                (error) => {
                    this.setState({
                        tableError: error,
                    });
                    alert(error.message); // TODO not firing
                }
            );
    }
    componentDidMount() {
        this.loadTabs();
        this.loadTable();
        this.loadColumns();
    }

    componentDidUpdate(prevProps, prevState) {
        if (prevState["tableName"] !== this.state.tableName) {
            this.loadTable();
            this.loadColumns();
        }
    }
    // method for handling when clicks on "Add New ..." button
    addRow() {
        if (!this.state.columnsLoaded || this.state.columnsError != null) {
            return;
        }
        const emptyRow = {};

        if (this.state.isEditMode) {
            const columnData = this.state.columns;

            const i = columnData.length + 1;
            emptyRow["Column_name"] = "col_" + i;
            emptyRow["Column_type"] = "varchar(45)";
            emptyRow["Referenced_table_name"] = "";

            this.setState({
                columns: [emptyRow, ...columnData]
            });
        } else {
            this.state.columns.map((v) => {
                const col = capitalizeFirstLetter(v.Column_name || "Table");
                emptyRow[col] = '';
            });
            const tableData = this.state.tableData;
            this.setState({
                tableData: [emptyRow, ...tableData]
            });
        }
    }
    // edit table handler
    editTable() {
        this.setState({
            isEditMode: !this.state.isEditMode
        });
    }
    // open the delete modal
    openModal() {
        this.setState({
            modalOpen: true
        });
    }
    // close the delete modal
    closeModal() {
        this.setState({
            modalOpen: false
        });
    }
    // get the id of the row that the user wishes to remove
    selectRow(rowId) {
        this.setState({
            selectedRowId: rowId
        })
    }

    render() {
        return (
            <div className="App">
                <Header />

                <div className="container">
                    <h1>{capitalizeFirstLetter(this.state.tableName)}</h1>

                    <Navbar
                        onClick={(tableName) => this.changeTable(tableName)}
                        items={this.state.tabs}
                        isLoaded={this.state.tabsLoaded}
                        activeTable={this.state.tableName}
                        error={this.state.tabsError}/>
                    <Buttons
                        tableName={this.state.tableName}
                        addRow={() => this.addRow()}
                        editTable={() => this.editTable()}
                        isEditMode={this.state.isEditMode}/>
                    <Table
                        items={this.state.isEditMode ? this.state.columns : this.state.tableData}
                        isLoaded={this.state.tableLoaded}
                        error={this.state.tableError}
                        tableName={this.state.tableName}
                        openModal={() => this.openModal()}
                        closeCellHandler={this.closeCellHandler}
                        selectRow={this.selectRow}
                        isEditMode={this.state.isEditMode}/>
                </div>
                <DeleteModal
                    handleClose={() => this.closeModal()}
                    handleDelete={() => this.handleDelete()}
                    show={this.state.modalOpen}
                    tableName={this.state.tableName}/>
            </div>
        );
    }
}
// the header
const Header = () => {
    return <header className="App-header">
                <h1>Welcome to Tabula Rasa</h1>
            </header>
};
// the "Add New ..." button component
const Add = (props) => {
    let itemName = props.tableName === 'Tables' ? 'Table' : props.tableName;
    if (props.isEditMode) itemName = "Column";
    const btnTxt = "Add New " + itemName;

    return <div className="col-md-1" style={{marginRight: btnTxt.length/2.5 + '%'}}>
                <button onClick={props.onClick}
                   className="btn btn-success">{btnTxt}</button>
            </div>
};
// the "Edit Table ..." button component
const Edit = (props) => {
    if (props.tableName === 'Tables') {
        return null;
    }
    const btnTxt = props.isEditMode ? "Exit Edit Mode" : "Edit Table " + props.tableName;
    const className = props.isEditMode ? "btn btn-default" : "btn btn-warning";

    return <div className="col-md-1">
                <button onClick={props.onClick}
                   className={className}>{btnTxt}</button>
            </div>
};
const Buttons = (props) => {
    return <div className="row"
                style={{marginBottom: '30px'}}>
                <Add onClick={props.addRow}
                     tableName={capitalizeFirstLetter(props.tableName)}
                     isEditMode={props.isEditMode}/>
                <Edit onClick={props.editTable}
                      tableName={capitalizeFirstLetter(props.tableName)}
                      isEditMode={props.isEditMode}/>
            </div>
};
// a tab component
const Tab = (props) => {
    return (
        <li role="presentation"
            onClick={() => props.onClick(props.tableName)}
            className={props.className}><a>{capitalizeFirstLetter(props.tableName)}</a></li>
    );
};
// the component with all the tabs
const Tabs = (props) => {
    const { error, isLoaded, items, onClick, activeTable } = props;
    if (error) {
        return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
        return <div>Loading...</div>;
    } else {
        return (items.map((item, i) => {
                const tableName = items[i].Table_name;
                return <Tab
                        onClick={onClick}
                        className={tableName === activeTable ? 'active' : ''}
                        tableName={tableName}/>
            })
        );
    }
};
// the Navbar component
const Navbar = (props) => {
    return (
        <div className="navbar">
            <div id="tabs">
                <ul className="nav nav-tabs">
                    <Tabs
                        error={props.error}
                        isLoaded={props.isLoaded}
                        items={props.items}
                        onClick={props.onClick}
                        activeTable={props.activeTable}/>
                </ul>
            </div>
        </div>
    );
};
// the row Delete button component
const RowDeleteButton = (props) => {
    return <td className="col-md-1">
            <button
                style={{padding: '0 10px 0 10px'}}
                role="button"
                data-toggle="modal"
                data-target="#confirm-delete"
                className="btn btn-danger"
                onClick={props.onClick}>Delete</button></td>
};
// the table head component
const THead = (props) => {
    if (!props.isLoaded) {
        return;
    }
    const keys = Object.keys(props.items[0]);

    return <thead>
    <tr>
        {
            keys.map((k, i) => {
                return <th key={i}>{k}</th>
            })
        }
        <th>Action</th>
    </tr>
    </thead>
};
// the table cell component
class Td extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isEdit: this.props.isEditMode,
        };
    }

    openCell() {
        this.setState({
            isEdit: true
        });
    }

    closeCellHandler(e) {
        this.props.closeCellHandler(e);
        this.setState({
            isEdit: false
        });
    }

    componentDidUpdate() {
        if (this.state.isEdit ||
            (this.props.isEdit && this.props.columnName === "Column_name")) {
           this.cellInput.focus();
        }
    }

    handleKeyPress(e) {
        if (e.key === 'Enter') {
            this.closeCellHandler(e);
        }
    }

    render() {
        if (this.props.columnName !== 'id' && this.state.isEdit || this.props.isEditMode) {
            return <td className={this.props.className}>
                        <input name={this.props.columnName}
                               defaultValue={this.props.cellData}
                               rowIndex={this.props.rowIndex}
                               placeholder={this.props.columnName}
                               style={{textAlign: 'center'}}
                               onBlur={(e) => this.closeCellHandler(e)}
                               ref={(input) => this.cellInput = input}
                               onKeyPress={(e) => this.handleKeyPress(e)}
                               type="text"
                               key={this.props.cellData}/></td>
        }
        const onClick = this.props.columnName === 'Id' ? null : () => this.openCell();

        return <td onClick={onClick}
                   className={this.props.className}>{this.props.cellData}</td>;
    }
}
// the table row component
const Tr = (props) => {
    let keys = Object.keys(props.rowData);

    const handleClick = () => {
        props.openModal();
        props.selectRow(props.rowData.Id || props.rowData.Table || props.rowData.Column_name);
    };

    return <tr>
        {
            keys.map((k, i) => {
                return <Td className={i === 0 ? 'col-md-1' : 'col-md-2'}
                           cellData={props.rowData[k]}
                           columnName={k}
                           rowIndex={props.rowIndex}
                           closeCellHandler={props.closeCellHandler}
                           isEditMode={props.isEditMode}/>
            })
        }
        <RowDeleteButton onClick={handleClick}/>
    </tr>
};
// the table body component
const TBody = (props) => {
    return <tbody>
        {
            props.items.map((item, i) => {
                return <Tr rowData={item}
                           openModal={props.openModal}
                           selectRow={props.selectRow}
                           rowIndex={i}
                           closeCellHandler={props.closeCellHandler}
                           isEditMode={props.isEditMode}/>
            })
        }
    </tbody>;
};
// the table component
const Table = (props) => {
    const { error, isLoaded, items, tableName, openModal, selectRow, isEditMode } = props;
    if (error) {
        return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
        return <div>Loading...</div>;
    } else if (!items || !items.length) {
        return null;
    } else {
        return (
            <div className="row">
                <div className="col-md-12">
                    <table className="table table-bordered">
                        <THead isLoaded={isLoaded}
                               items={items}
                               isEditMode={props.isEditMode}/>
                        <TBody items={items}
                               tableName={tableName}
                               openModal={openModal}
                               selectRow={selectRow}
                               closeCellHandler={props.closeCellHandler}
                               isEditMode={props.isEditMode}/>
                    </table>
                </div>
            </div>
        );
    }
};
// the delete modal component
const DeleteModal = (props) => {
    return (
        <div>
            <Modal show={props.show}
                   onHide={props.handleClose}>
                <Modal.Header closeButton>
                    <Modal.Title>Delete {props.tableName === 'tables' ? 'Table'
                        : capitalizeFirstLetter(props.tableName)}</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <h5>Are you sure you want to delete the item from your collection?</h5>
                </Modal.Body>
                <Modal.Footer>
                    <Button onClick={props.handleClose}>Close</Button>
                    <Button className="btn btn-danger"
                            onClick={() => props.handleDelete()}>Delete</Button>
                </Modal.Footer>
            </Modal>
        </div>
    );
};
// helper function to capitalize the first letter in a string
const capitalizeFirstLetter = (string) => {
    return string.charAt(0).toUpperCase() + string.slice(1);
};

ReactDOM.render(<App />, document.getElementById('root'));
//registerServiceWorker();