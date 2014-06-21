/**
 * @jsx React.DOM
 */

var Skit = React.createClass({
  handleDelete: function() {
    $.get('/s/'+this.props.hash+'/delete', function(data) {
      window.State.refresh();
    });
    return false;
  },
  handleClick: function() {
    window.State.goto('/s/'+this.props.hash);
    return false;
  },
  render: function() {
    return (
      <div className="ui-item" style={{"border-color": handleColor(this.props.user)}}>
        <a href="#" onClick={this.handleClick}>{this.props.children}</a>
        <a className="ui-item-delete" href="#" onClick={this.handleDelete}>x</a>
      </div>
    );
  }
});

var SkitList = React.createClass({
  render: function() {
    if (!this.props.data) return (<div></div>);
    var skitNodes = this.props.data.map(function(skit) {
      return <Skit key={skit.hash} hash={skit.hash} user={skit.user} root={skit.root}>{skit.text}</Skit>;
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

var SkitHeader = React.createClass({
  handleBack: function(e) {
    window.State.back();
    return false;
  },
  render: function() {
    if (!this.props.data) {
      return (
        <header>
          <h1>Skit</h1>
        </header>
      )
    }
    return (
      <header>
        <a className="ui-back" href="#" onClick={this.handleBack}>Back</a>
        <h1>Skit {this.props.data.text}</h1>
      </header>
    )
  }
});

var SkitBox = React.createClass({
  handleSubmit: function(skit) {
    var skits = this.state.data;
    var newSkits = [skit].concat(skits);

    this.setState({data: newSkits});

    $.post('/s/save', skit, function(data) {
      window.State.refresh();
    }.bind(this));
    return false;
  },
  handleMessage: function(e) {
    var data = JSON.parse(e.data)
    this.setState({data: data});
  },
  componentWillMount: function() {
    var url = this.props.url;
    window.SOCKET.onmessage = this.handleMessage;
    window.SOCKET.onopen = function() {
      window.State.goto(url);
    }
  },
  getInitialState: function() {
    return {data: []};
  },
  render: function() {
    return (
      <div>
        <SkitHeader data={this.state.data.skit} />
        <article>
          <SkitForm onSubmit={this.handleSubmit} parent={this.state.data.skit} />
          <SkitList data={this.state.data.children} />
        </article>
      </div>
    );
  }
});
