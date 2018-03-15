module.exports = {
    "extends": "airbnb",
    "plugins": ["jest"],
    "env": {
        "browser": true,
        "jest/globals": true
    },
    "rules": {
        "import/no-extraneous-dependencies": [
            "error", {
                "devDependencies": true,
                "optionalDependencies": false,
                "peerDependencies": false
            }
        ]
    }
};
