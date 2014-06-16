/**
 * @jsx React.DOM
 */

var Skit = React.createClass({
  handleDelete: function() {
    $.get('/skit/delete/'+this.props.hash, function(data) {
      window.SOCKET.send('{"url": "/"}');
    });
  },
  handleColor: function() {
    var hue_resolution = 359;
    var hue = parseInt(this.props.user, 16) % hue_resolution;
    return 'hsl(' + hue + ', '
                            + (.7 * 100) + '%, '
                            + (.7 * 100)+'%)';
  },
  render: function() {
    return (
      <div className="ui-item" style={{"border-color": this.handleColor()}}>
        <a href={'/skit/view/'+this.props.hash}>{this.props.children}</a>
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

var SkitBox = React.createClass({
  handleSubmit: function(skit) {
    var skits = this.state.data;
    var newSkits = [skit].concat(skits);
    this.setState({data: newSkits});

    $.post('/skit/save/', skit, function(data) {
      window.SOCKET.send('{"url": "/"}');
    }.bind(this));
  },
  handleMessage: function(e) {
    var data = JSON.parse(e.data)
    this.setState({data: data.skits});
  },
  componentWillMount: function() {
    window.SOCKET.onmessage = this.handleMessage;
    window.SOCKET.onopen = function() {
      window.SOCKET.send('{"url": "/"}');
    }
  },
  getInitialState: function() {
    return {data: []};
  },
  render: function() {
    return (
      <div>
        <header>
          <h1>Skits</h1>
        </header>
        <article>
          <SkitForm onSubmit={this.handleSubmit} />
          <SkitList data={this.state.data} />
        </article>
      </div>
    );
  }
});
