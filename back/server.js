const express = require('express')
const jwt = require('jsonwebtoken')
const uuid = require('uuid')
const app = express()
const port = 3000

let json =
{
    "type": "collection", "count": 3, "commandes": [
        {
            "suuid": "JSESSIONID=86F633899E8C2CF1577717",
            "cameraCode": "11010000001320000050%2316%234c6a6bade5c74d66b",
            "msg": "rtsp://stream.strba.sk:1935/strba/VYHLAD_JAZERO.stream",
            "code": 200
        },
        {
            "suuid": "JSESSIONID=86F633899E8C2CF1577718",
            "cameraCode": "11010000001320000050%2316%234c6a6bade5c74d67b",
            "msg": "rtsp://stream.strba.sk:1935/strba/VYHLAD_JAZERO.stream",
            "code": 200
        },
        {
            "suuid": "JSESSIONID=86F633899E8C2CF1577719",
            "cameraCode": "11010000001320000050%2316%234c6a6bade5c74d68b",
            "msg": "rtsp://stream.strba.sk:1935/strba/VYHLAD_JAZERO.stream",
            "code": 200
        },
    ],
}

app.use(express.json()) // for parsing application/json
app.use(express.urlencoded({ extended: true })) // for parsing application/x-www-form-urlencoded

app.get('/', (req, res) => {
    res.send('Hello Ziyi! Have a nice day !')
})

app.get("/videos", auth, async(req, res, next) => {
    let commandes = []
    json.commandes.forEach(commande => {
        commandes.push({
            suuid: commande.suuid,
            cameraCode: commande.cameraCode,
            msg: commande.msg,
            code: commande.code

        });
    });
    res.json(commandes)
})

app.post('/login', function (req, res){
    const username = req.body.username
    const user = {name: username}

    // sign with RSA SHA256
    const accessToken = jwt.sign({user}, 'my_secret_key', {algorithm: 'HS256', expiresIn: '600s'})
    return res.status(201).json({
        user: req.headers.nom,
        token: accessToken
    });
})

function auth(req, res, next) {
    if(req.headers['authorization']){
        const authHeader = req.headers['authorization']
        const token = authHeader && authHeader.split(' ')[1]
        if(token !== null){
            jwt.verify(token, 'my_secret_key', (err, user) => {
                if(err){
                    res.status(500).json({
                        error: "Not Authorized"
                    });
                    throw new Error("Not Authorized");
                }
                req.user = user
                next();
            })
        }
    }else{
        return res.status(401).json({
            error: "Missing credential"
        })
    }
}

app.use((function (req, res) {
    res.sendStatus(404)
}))

app.listen(port, () => {
    console.log(`Example app listening at http://localhost:${port}`)
})




