const {Button, Modal} = ReactBootstrap;

class App extends React.Component {
    constructor(props) {
        super(props);
        this.selectRow = this.selectRow.bind(this);
        this.closeCellHandler = this.closeCellHandler.bind(this);
        this.state = {
            tableName: 'home',
            modalOpen: false
        }
    }

    changeTable(nextTable) {
        this.setState({
            tableName: nextTable,
        });
    }

    loadTabs() {
        fetch("/rest/home")
            .then(res => res.json(), {
                method: 'GET'
            })
            .then(
                (result) => {
                    result = [{table_name: 'home'}].concat(result);
                    this.setState({
                        tabsLoaded: true,
                        tabs: result
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

    loadTable() {
        fetch("/rest/" + this.state.tableName, {
            method: 'GET'
        })
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        tableLoaded: true,
                        tableData: result
                    });
                },
                (error) => {
                    this.setState({
                        tableLoaded: true,
                        tableError: error,
                    });
                }
            )
    }

    getColumns() {
        fetch("/rest/" + this.state.tableName + "/cols", {
            method: 'GET'
        })
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        columns: result,
                        columnsLoaded: true
                    });
                },
                (error) => {
                    this.setState({
                        columnsLoaded: true,
                        columnsError: error,
                    });
                }
            )
    }

    saveTable(data, afterSave) {
        fetch("/rest/" + this.state.tableName, {
            method: 'POST',
            body: JSON.stringify(data)
        })
            .then(res => res.json())
            .then(result => {
                afterSave(result);
            });
    }

    closeCellHandler(event) {
        const tableData = this.state.tableData.slice(),
            {value, name, attributes} = event.target,
            index = Number(attributes.rowindex.value);

        tableData[index][name] = value;

        this.saveTable(tableData[index], (result) => {
            if (!tableData[index].id) {
                tableData[index].id = result.id;
            }
            this.setState({
                tableData: tableData,
                tableLoaded: true
            });
        });
    }

    componentDidMount() {
        this.loadTabs();
        this.loadTable();
        this.getColumns();
    }

    componentDidUpdate(prevProps, prevState) {
        if (prevState["tableName"] !== this.state.tableName) {
            this.loadTable();
            this.getColumns();
        }
    }

    addRow() {
        if (!this.state.columnsLoaded || this.state.columnsError != null) {
            return;
        }
        const emptyRow = {};
        this.state.columns.map((v) => {
            const col = v.column_name;
            emptyRow[col] = '';
        });
        const tableData = this.state.tableData;
        this.setState({
            tableData: [emptyRow, ...tableData]
        });
    }

    handleDeleteRow() {
        const tableName = this.state.tableName;
        const id = this.state.selectedRowId;

        fetch("/rest/" + tableName + "?id=" + id, {
            method: 'DELETE'
        })
            .then(
                (result) => {
                    this.setState({
                        tableData: this.state.tableData.filter((v) => v.id !== id ),
                        modalOpen: false,
                        selectedRowId: null
                    });
                },
                (error) => {
                    /*this.setState({
                        tableError: error,
                    });*/
                }
            );
    }

    openModal() {
        this.setState({
            modalOpen: true
        });
    }

    closeModal() {
        this.setState({
            modalOpen: false
        });
    }

    selectRow(rowId) {
        this.setState({
            selectedRowId: rowId
        })
    }

    render() {
        return (
            <div className="App">
                <header className="App-header">
                    <img src='/static/src/logo.svg'
                         className="App-logo"
                         alt="logo"/>
                    <h1 className="App-title">Welcome to Tabula Rasa</h1>
                </header>

                <div className="container">
                    <h1>{capitalizeFirstLetter(this.state.tableName)}</h1>

                    <Navbar
                        onClick={(tableName) => this.changeTable(tableName)}
                        items={this.state.tabs}
                        isLoaded={this.state.tabsLoaded}
                        activeTable={this.state.tableName}
                        error={this.state.tabsError}/>
                    <div className="row"
                         style={{marginBottom: '30px'}}>
                        <div className="col-md-1">
                            <Add onClick={() => {this.addRow()}}
                                 tableName={capitalizeFirstLetter(this.state.tableName)}/>
                        </div>
                    </div>
                    <div className="row">
                        <div className="col-md-12">
                            <Table
                                items={this.state.tableData}
                                isLoaded={this.state.tableLoaded}
                                error={this.state.tableError}
                                tableName={this.state.tableName}
                                openModal={() => this.openModal()}
                                closeCellHandler={this.closeCellHandler}
                                selectRow={this.selectRow}/>
                        </div>
                    </div>
                </div>
                <DeleteModal
                    handleClose={() => this.closeModal()}
                    handleDelete={() => this.handleDeleteRow()}
                    show={this.state.modalOpen}
                    tableName={this.state.tableName}/>
            </div>
        );
    }
}

const Add = (props) => {
    const tableName = props.tableName === 'Home' ? 'Table' : props.tableName;
    return <button onClick={props.onClick}
                   className="btn btn-default">Add New {tableName}</button>
};

const Tab = (props) => {
    return (
        <li role="presentation"
            onClick={() => props.onClick(props.tableName)}
            className={props.className}><a>{capitalizeFirstLetter(props.tableName)}</a></li>
    );
};

const Tabs = (props) => {
    const { error, isLoaded, items, onClick, activeTable } = props;
    if (error) {
        return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
        return <div>Loading...</div>;
    } else {
        return (items.map((item, i) => {
                const tableName = items[i].table_name;
                return <Tab
                        onClick={onClick}
                        className={tableName === activeTable ? 'active' : ''}
                        tableName={tableName}/>
            })
        );
    }
};

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

const RowDeleteButton = (props) => {
    return <td className="col-md-1">
            <button
                role="button"
                data-toggle="modal"
                data-target="#confirm-delete"
                className="btn btn-danger"
                onClick={props.onClick}>Delete</button></td>
};

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

class Td extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isEdit: false,
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
        if (this.state.isEdit) {
           this.cellInput.focus();
        }
    }

    render() {
        if (this.props.columnName !== 'id' && this.state.isEdit) {
            return <td className={this.props.className}>
                        <input name={this.props.columnName}
                               defaultValue={this.props.cellData}
                               rowIndex={this.props.rowIndex}
                               placeholder={this.props.columnName}
                               style={{textAlign: 'center'}}
                               onBlur={(e) => this.closeCellHandler(e)}
                               ref={(input) => { this.cellInput = input; }}
                               type="text"/></td>
        }
        return <td onClick={() => this.openCell()}
                   className={this.props.className}>{this.props.cellData}</td>;
    }
}

const Tr = (props) => {
    let keys = Object.keys(props.rowData);

    const handleClick = () => {
        props.openModal();
        props.selectRow(props.rowData.id)
    };

    return <tr>
        {
            keys.map((k, i) => {
                return <Td className={i === 0 ? 'col-md-1' : 'col-md-2'}
                           cellData={props.rowData[k]}
                           columnName={k}
                           rowIndex={props.rowIndex}
                           closeCellHandler={props.closeCellHandler} />
            })
        }
        <RowDeleteButton onClick={handleClick}/>
    </tr>
};

const TBody = (props) => {
    return <tbody>
        {
            props.items.map((item, i) => {
                return <Tr rowData={item}
                           openModal={props.openModal}
                           selectRow={props.selectRow}
                           rowIndex={i}
                           closeCellHandler={props.closeCellHandler}/>
            })
        }
    </tbody>;
};

const Table = (props) => {
    const { error, isLoaded, items, tableName, openModal, selectRow } = props;
    if (error) {
        return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
        return <div>Loading...</div>;
    } else if (!items.length) {
        return null;
    } else {
        return (
            <table className="table table-bordered">
                <THead isLoaded={isLoaded}
                       items={items}/>
                <TBody items={items}
                       tableName={tableName}
                       openModal={openModal}
                       selectRow={selectRow}
                       closeCellHandler={props.closeCellHandler}/>
            </table>
        );
    }
};

const DeleteModal = (props) => {
    return (
        <div>
            <Modal show={props.show}
                   onHide={props.handleClose}>
                <Modal.Header closeButton>
                    <Modal.Title>Delete {props.tableName === 'home' ? 'Table'
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

const capitalizeFirstLetter = (string) => {
    return string.charAt(0).toUpperCase() + string.slice(1);
};

ReactDOM.render(<App />, document.getElementById('root'));
//registerServiceWorker();