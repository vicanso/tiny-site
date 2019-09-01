import React from "react";
import PropTypes from "prop-types";

class FileZone extends React.Component {
    render() {
        return <div className="FileZone">

        </div>
    }
}

FileZone.propTypes = {
    account: PropTypes.string.isRequired,
    roles: PropTypes.array.isRequired
};

export default FileZone;