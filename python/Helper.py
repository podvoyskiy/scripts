import subprocess

class Helper:
    
    @staticmethod
    def exec_cli_cmd(cmd: str) -> str:
        proc = subprocess.Popen(cmd, stdout=subprocess.PIPE, shell=True)
        (stdout, stderr) = proc.communicate()

        if (stderr is not None): return stderr.decode("utf-8")
        return stdout.decode("utf-8")