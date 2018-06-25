'use strict';
var dataProvider = require('../data/workHistory.js');
/**
 * Operations on /workHistory
 */
module.exports = {
    /**
     * summary: Retrieves work history.
     * description: 
     * parameters: 
     * produces: application/json, text/json
     * responses: 200
     */
    get: function workHistory_get(req, res, next) {
        /**
         * Get the data for response 200
         * For response `default` status 200 is used.
         */
        var status = 200;
        var provider = dataProvider['get']['200'];
        provider(req, res, function (err, data) {
            if (err) {
                next(err);
                return;
            }
            res.status(status).send(data && data.responses);
        });
    }
};
