class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            tableName: 'home',
            deleteModal: {
                active: false,
                title: 'U sure u wanna delete?'
            }
        }
    }

    changeTable(nextTable) {
        this.setState({
            tableName: nextTable,
        });
    }

    loadTabs() {
        fetch("/rest/home")
            .then(res => res.json())
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
        fetch("/rest/" + this.state.tableName)
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

    componentDidMount() {
        this.loadTabs();
        this.loadTable();
    }

    componentDidUpdate(prevProps, prevState) {
        if (prevState["tableName"] !== this.state.tableName) {
            this.loadTable();
        }
    }

    hideDeleteModal() {
        this.setState({
            deleteModal: {
                active: !this.state.deleteModal.active
            }
        });
    }

    render() {
        return (
            <div className="App">
                <header className="App-header">
                    <img src='/static/src/logo.svg' className="App-logo" alt="logo" />
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

                    <div className="row" style={{marginBottom: '30px'}}>
                        <div className="col-md-12">
                            <Table
                                items={this.state.tableData}
                                isLoaded={this.state.tableLoaded}
                                error={this.state.tableError}
                                />
                        </div>
                    </div>
                </div>
                <DeleteModal onClickCancel={() => {this.hideDeleteModal()}}
                             onConfirm={() => {this.hideDeleteModal()}}
                             isActive={this.state.deleteModal.active}
                             title={this.state.deleteModal.title} />
            </div>
        );
    }
}

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

const THead = (props) => {
    if (!props.isLoaded) {
        return;
    }
    let keys = Object.keys(props.items[0]);

    return <thead>
    <tr>
        {
            keys.map((k, i) => {
                return <th key={i}>{k}</th>
            })
        }
    </tr>
    </thead>
};

class Td extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isEdit: false
        }
    }

    openCell() {
        this.setState({
            isEdit: true
        });
    }

    closeCell() {
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
        if (this.state.isEdit) {
            return <td className={this.props.className}><input onBlur={() => this.closeCell()}
                                                               ref={(input) => { this.cellInput = input; }}  type="text"/></td>
        }
        return <td onClick={() => this.openCell()} className={this.props.className}>{this.props.cellData}</td>;
    }
}

const Tr = (props) => {
    let keys = Object.keys(props.rowData);

    return <tr>
        {
            keys.map((k, i) => {
                return <Td className={i === 0 ? 'col-md-1' : 'col-md-2'} cellData={props.rowData[k]}/>
            })
        }
    </tr>
};

const TBody = (props) => {
    return <tbody>
        {
            props.items.map((item, i) => {
                return <Tr rowData={item} />
            })
        }
    </tbody>;
};

const Table = (props) => {
    const { error, isLoaded, items } = props;
    if (error) {
        return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
        return <div>Loading...</div>;
    } else {
        return (
            <table className="table table-bordered">
                <THead isLoaded={isLoaded} items={items}/>
                <TBody items={items}/>
            </table>
        );
    }
};

const DeleteModal = (props) => {
    return (
        <div className={
                props.isActive ? 'modal' : 'modal fade'}
            style={{display:
                props.isActive ? 'inline' : 'none'}}
            id="confirm-delete" tabIndex="-1" role="dialog"
             aria-labelledby="myModalLabel">
            <div className="modal-dialog">
                <div className="modal-content">
                    <div className="modal-header title">{props.title}</div>
                    <div className="modal-footer">
                        <button onClick={props.onClickCancel} type="button" className="btn btn-default">Cancel</button>
                        <button onClick={props.onConfirm} type="button" className="btn btn-danger btn-ok">Delete</button>
                    </div>
                </div>
            </div>
        </div>
    );
};

const capitalizeFirstLetter = (string) => {
    return string.charAt(0).toUpperCase() + string.slice(1);
};

ReactDOM.render(<App />, document.getElementById('root'));
//registerServiceWorker();
