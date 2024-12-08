import asyncio
from aiosmtpd.controller import Controller
from email import message_from_string

class CustomMessageHandler:
    async def handle_DATA(self, server, session, envelope):
        peer = session.peer
        mailfrom = envelope.mail_from
        rcpttos = envelope.rcpt_tos
        data = envelope.content.decode('utf-8', errors='replace')

        print(f"Connection from: {peer}")
        print(f"From: {mailfrom}")
        print(f"To: {rcpttos}")
        print(f"Message:\n{data}")
        print("----- End of Message -----\n")

        # 如果需要存储或进一步处理，可以在此添加逻辑
        message = message_from_string(data)
        print(message["Subject"])


        # 返回接受状态
        return '250 Message accepted for delivery'

if __name__ == "__main__":
    handler = CustomMessageHandler()

    # 创建 SMTP 服务器
    controller = Controller(handler, hostname='0.0.0.0', port=25)
    print("SMTP server is running on 0.0.0.0:25")

    # 启动服务器
    controller.start()

    try:
        asyncio.get_event_loop().run_forever()
    except KeyboardInterrupt:
        print("Server stopped.")
        controller.stop()
