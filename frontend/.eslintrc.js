module.exports = {
    "extends": "airbnb",
    "parser": "babel-eslint",
    "plugins": ["jest", "babel"],
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
