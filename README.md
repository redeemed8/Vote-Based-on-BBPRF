# ***Vote-Based-on-VOPRF***

## 基于Vote-Based-on-VOPRF和SM2的投票系统



## ***1.服务器密钥生成***

$$
\begin{equation}
	公钥pk_s\quad [ N,\quad G,\quad Y,\quad H,\quad ct_u,\quad ct_y]
\end{equation}
$$

$$
\begin{equation}
私钥sk_s\quad[x,\quad u,\quad y]
\end{equation}
$$

## ***2.客户端注册并获取公钥***

##### **·** &nbsp; 使用登录token到认证服务换取对应的私钥，并获取私钥的签名u<sub>c</sub>

$$
\begin{equation}
私钥sk_c\quad [u_c,\quad uc_{sign},\quad a,\quad b,\quad r]
\end{equation}
$$

## ***3.客户端请求给加密token签名***

$$
\begin{equation}
digest\quad {\small ←}\quad ct_β \quad {\small ←}\quad Enc(pk_{cs},am+bp)+a*ct_u+a*r*ct_y
\end{equation}
$$

$$
\begin{equation}
state\quad {\small ←}\quad (a,b)
\end{equation}
$$

## ***4.服务端给token签名***

$$
\begin{equation}
(u,e)\quad {\small ←}\quad ct_β \quad {\small ←}\quad digest
\end{equation}
$$

$$
\begin{equation}
β\quad {\small ←}\quad Dec(sk_{cs},ct_β)\quad {\small ←}\quad ((e/u^x)-1\ mod\ n)/N
\end{equation}
$$













