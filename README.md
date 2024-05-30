# ***Vote-Based-on-VOPRF***

## 基于Vote-Based-on-VOPRF和SM2的投票系统



## ***1.服务器密钥生成***

$$
\begin{equation}
	公钥pk_s \quad {\small ←}\quad pk_{VOPRF} \quad {\small ←}\quad [\ N,\quad G,\quad Y,\quad H,\quad ct_u,\quad ct_y\ ]
\end{equation}
$$

$$
\begin{equation}
私钥sk_s \quad {\small ←}\quad sk_{VOPRF} \quad {\small ←}\quad [\ x,\quad u,\quad y\ ]
\end{equation}
$$



## ***2.客户端注册并获取公钥***

##### **·** &nbsp; 使用登录token到认证服务换取对应的私钥，并获取私钥的签名u<sub>c</sub>

$$
\begin{equation}
私钥sk_c\quad [\ u_c,\quad uc_{sign},\quad a,\quad b,\quad r \ ]
\end{equation}
$$



## ***3.客户端请求给盲化token签名***

$$
\begin{equation}
digest\quad {\small ←}\quad ct_β \quad {\small ←}\quad Enc(pk_{cs}, \  am+bp) + a * ct_u + a * r * ct_y
\end{equation}
$$

$$
\begin{equation}
state \quad {\small ←}\quad (a, \  b)
\end{equation}
$$

$$
\begin{equation}
tag \quad {\small ←}\quad v \quad {\small ←}\quad G^{[{u_c+H(mag)]}^{-1}} \quad {\small ←}\quad F_{DY}(u_c,\ H(msg))
\end{equation}
$$

## ***4.服务端给token签名***

##### **·** &nbsp; 将&nbsp;ct<sub>β</sub>&nbsp;转换成二元组

$$
\begin{equation}
(u,\ e)\quad {\small ←}\quad ct_β \quad {\small ←}\quad digest
\end{equation}
$$

##### **·** &nbsp; 计算出对应的β并得到盲化伪随机函数

$$
\begin{equation}
β\quad {\small ←}\quad [\ (e/u^x)-1\ mod\ n\ ]\ /\ N \quad {\small ←}\quad Dec(\ sk_{cs},\ ct_β\ )  \quad {\small ←}\quad Eval(\ sk_{VOPRF},\ digest\ )
\end{equation}
$$

$$
\begin{equation}
blindPRF \quad {\small ←}\quad \{ \ F =G^{β^{-1}}, \ ... \  \} 
\end{equation}
$$



## ***5.客户端对签名token进行去盲化***

$$
\begin{equation}
tok  \quad {\small ←}\quad \tau \quad {\small ←}\quad F^a \quad {\small ←}\quad Decode(\ state,blindPRF\ )
\end{equation}
$$



## ***6.向服务器验证token***

$$
\begin{equation}
verify \quad {\small ←}\quad (u_c,\ uc_{sign},\ msg,\ token,\ pk_s,\ sk_c,\ tag)
\end{equation}
$$

$$
\begin{equation}
true\quad {\small ←}\quad token\ {\small ==}\ G^{{[m+u+ry]}^{-1}} \quad {\small ←}\quad tag\ not\ exist \quad {\small ←}\quad verify
\end{equation}
$$







