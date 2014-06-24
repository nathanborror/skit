/**
 * @jsx React.DOM
 */

var Skit = React.createClass({
  handleDelete: function() {
    $.get('/s/'+this.props.hash+'/delete');
    this.setState({visible: false});
    return false;
  },
  handleClick: function() {
    this.props.pushSkitBox('/s/'+this.props.hash);
    return false;
  },
  getInitialState: function() {
    return {visible: true};
  },
  render: function() {
    return this.state.visible ? (
      <div className="ui-item" style={{"border-color": handleColor(this.props.user)}}>
        <a href="#" onClick={this.handleClick}>{this.props.children}</a>
        <a className="ui-item-delete" href="#" onClick={this.handleDelete}>x</a>
      </div>
    ) : (<span />);
  }
});

var SkitList = React.createClass({
  render: function() {
    if (!this.props.data) return (<div></div>);
    var pushSkitBox = this.props.pushSkitBox;
    var skitNodes = this.props.data.map(function(skit) {
      return <Skit
        key={skit.hash}
        hash={skit.hash}
        user={skit.user}
        pushSkitBox={pushSkitBox}
        root={skit.root}>{skit.text}</Skit>;
    });
    return (
      <div className="ui-list">
        {skitNodes}
      </div>
    );
  }
});

var SkitForm = React.createClass({
  handleSubmit: function() {
    var hash = this.refs.hash.getDOMNode().value.trim();
    var parent = this.refs.parent.getDOMNode().value.trim();
    var root = this.refs.root.getDOMNode().value.trim();
    var text = this.refs.text.getDOMNode().value.trim();

    this.props.onSubmit({hash: hash, parent: parent, root: root, text: text});
    this.refs.hash.getDOMNode().value = '';
    this.refs.text.getDOMNode().value = '';
    return false;
  },
  render: function() {
    var hash = this.props.parent ? this.props.parent.hash : "";
    var root = this.props.parent ? this.props.parent.root : "";
    return (
      <form className="ui-form" onSubmit={this.handleSubmit}>
        <input type="text" ref="text" placeholder="What are you thinking about?" />
        <input type="hidden" ref="hash" value="" />
        <input type="hidden" ref="parent" value={hash} />
        <input type="hidden" ref="root" value={root} />
        <button type="submit" className="hidden">Save</button>
      </form>
    );
  }
});

var SkitBox = React.createClass({
  handleSubmit: function(skit) {
    // optimisticly add skit
    var current = this.state.data;
    current.children.unshift(skit);
    this.setState({data: current});

    $.post('/s/save', skit, function(data) {
      // replace skit with real skit
      current.children[0] = data.skit;
      this.setState({data: current});
    }.bind(this));
    return false;
  },
  handleMessage: function(data) {
    this.setState({data: data});
  },
  componentWillMount: function() {
    var url = this.props.url;
    window.SOCKET.subscribe(url, this.handleMessage);
    window.SOCKET.request(url);
    window.history.pushState({}, "", url);
  },
  getInitialState: function() {
    return {data: {}};
  },
  render: function() {
    return (
      <div>
        <article>
          <SkitForm onSubmit={this.handleSubmit} parent={this.state.data.skit} />
          <SkitList
            key={this.props.url}
            data={this.state.data.children}
            pushSkitBox={this.props.pushSkitBox} />
        </article>
      </div>
    );
  }
});

var SkitBoxes = React.createClass({
  pushSkitBox: function(url) {
    this.state.urls.push(url);
    this.setState(this.state);
  },
  popSkitBox: function() {
    if (this.state.urls.length == 1) {
      return;
    }
    this.state.urls.pop();
    this.setState(this.state);
  },
  getInitialState: function() {
    // delay setting state until the websocket is open
    var self = this;
    window.SOCKET.onopen = function() {
      self.setState({urls: self.props.urls});
    };
    return {urls: []};
  },
  render: function() {
    var root = true;
    var self = this;
    var boxes = this.state.urls.map(function(url) {
      var b = <SkitBox
        key={url}
        isRoot={root}
        url={url}
        pushSkitBox={self.pushSkitBox}
        popSkitBox={self.popSkitBox} />;
      root = false;
      return b;
    });
    // force to last item
    return (
      <div>{boxes}</div>
    );
  }
});
