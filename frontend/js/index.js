import React from 'react';
import ReactDOM from 'react-dom';

const App = () => {
    return (
        <div>
            <h1>Hello, React!</h1>
            <p>This is a simple React HTML page.</p>
        </div>
    );
};

ReactDOM.render(<App />, document.getElementById('root'));