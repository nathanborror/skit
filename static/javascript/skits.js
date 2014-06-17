/**
 * @jsx React.DOM
 */

var Skit = React.createClass({
  handleDelete: function() {
    $.get('/s/'+this.props.hash+'/delete', function(data) {
      window.SOCKET.send('{"url": "/"}');
    });
    return false;
  },
  render: function() {
    return (
      <div className="ui-item" style={{"border-color": handleColor(this.props.user)}}>
        <a href={'/s/'+this.props.hash}>{this.props.children}</a>
        <a className="ui-item-delete" href="#" onClick={this.handleDelete}>x</a>
      </div>
    );
  }
});

var SkitList = React.createClass({
  render: function() {
    if (!this.props.data) return (<div></div>);
    var skitNodes = this.props.data.map(function(skit) {
      return <Skit key={skit.hash} hash={skit.hash} user={skit.user}>{skit.text}</Skit>;
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
    var text = this.refs.text.getDOMNode().value.trim();

    this.props.onSubmit({hash: hash, parent: parent, text: text});
    this.refs.hash.getDOMNode().value = '';
    this.refs.text.getDOMNode().value = '';
    return false;
  },
  render: function() {
    return (
      <form className="ui-form" onSubmit={this.handleSubmit}>
        <input type="text" ref="text" placeholder="What are you thinking about?" />
        <input type="hidden" ref="hash" value="" />
        <input type="hidden" ref="parent" value={this.props.parent} />
        <button type="submit" className="hidden">Save</button>
      </form>
    );
  }
});

var SkitHeader = React.createClass({
  handleBack: function(e) {
    window.history.back();
    return false;
  },
  render: function() {
    if (!this.props.data) {
      return (
        <header>
          <h1>Skit</h1>
        </header>
      )
    } else {
      return (
        <header>
          <a className="ui-back" href="#" onClick={this.handleBack}>Back</a>
          <h1>Skit {this.props.data.text}</h1>
        </header>
      )
    }
  }
});

var SkitBox = React.createClass({
  handleSubmit: function(skit) {
    var skits = this.state.data;
    var newSkits = [skit].concat(skits);
    var payload = '{"url": "'+this.props.url+'"}';
    this.setState({data: newSkits});

    $.post('/s/save', skit, function(data) {
      window.SOCKET.send(payload);
    }.bind(this));
    return false;
  },
  handleMessage: function(e) {
    var data = JSON.parse(e.data)
    this.setState({data: data});
  },
  componentWillMount: function() {
    var payload = '{"url": "'+this.props.url+'"}';
    window.SOCKET.onmessage = this.handleMessage;
    window.SOCKET.onopen = function() {
      window.SOCKET.send(payload);
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
          <SkitForm onSubmit={this.handleSubmit} parent={this.props.parent} />
          <SkitList data={this.state.data.children} />
        </article>
      </div>
    );
  }
});
